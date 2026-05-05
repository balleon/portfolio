import pulumi
import json
import pulumi_aws as aws

vpc_cidr = "10.0.0.0/16"
private_subnet_cidrs = ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
public_subnet_cidrs = ["10.0.3.0/24", "10.0.4.0/24", "10.0.5.0/24"]
availability_zones = ["eu-west-3a", "eu-west-3b", "eu-west-3c"]

################################################################################
# VPC (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/vpc/)
################################################################################

vpc = aws.ec2.Vpc("vpc",
    cidr_block=vpc_cidr,
    tags={
        "Name": f"vpc-{pulumi.get_project()}",
        "environment": pulumi.get_stack(),
    }
)

################################################################################
# Private Subnets (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/subnet/)
################################################################################

private_subnets = []
for i in range(len(private_subnet_cidrs)):
    private_subnet = aws.ec2.Subnet(f"private_subnet_{availability_zones[i]}",
        vpc_id=vpc.id,
        cidr_block=private_subnet_cidrs[i],
        availability_zone=availability_zones[i],
        tags={
            "Name": f"private-subnet-{availability_zones[i]}",
            "environment": pulumi.get_stack(),
        }
    )
    private_subnets.append(private_subnet)

################################################################################
# Public Subnets (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/subnet/)
################################################################################

public_subnets = []
for i in range(len(public_subnet_cidrs)):
    public_subnet = aws.ec2.Subnet(f"public_subnet_{availability_zones[i]}",
        vpc_id=vpc.id,
        cidr_block=public_subnet_cidrs[i],
        availability_zone=availability_zones[i],
        tags={
            "Name": f"public-subnet-{availability_zones[i]}",
            "environment": pulumi.get_stack(),
        }
    )
    public_subnets.append(public_subnet)

################################################################################
# Internet Gateway (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/internetgateway/)
################################################################################

igw = aws.ec2.InternetGateway("igw",
    vpc_id=vpc.id,
    tags={
        "Name": f"igw-{pulumi.get_project()}",
        "environment": pulumi.get_stack(),
    }
)

################################################################################
# Route Table (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/routetable/)
################################################################################

public_rt = aws.ec2.RouteTable("public_rt",
    vpc_id=vpc.id,
    tags={
        "Name": "rt-public",
        "environment": pulumi.get_stack(),
    }
)

private_rt = aws.ec2.RouteTable("private_rt",
    vpc_id=vpc.id,
    tags={
        "Name": "rt-private",
        "environment": pulumi.get_stack(),
    }
)

################################################################################
# Elastic IP (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/eip/)
################################################################################

nat_gw_eip = aws.ec2.Eip("nat_gw_eip",
    tags={
        "Name": "nat-gw-eip",
        "environment": pulumi.get_stack(),
    }
)

################################################################################
# NAT Gateway (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/natgateway/)
################################################################################

nat_gw = aws.ec2.NatGateway("nat_gw",
    allocation_id=nat_gw_eip.id,
    subnet_id=public_subnets[0].id,
    tags={
        "Name": "nat-gw",
        "environment": pulumi.get_stack(),
    }
)

################################################################################
# Route (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/route/)
################################################################################

igw_route = aws.ec2.Route("igw_route",
    route_table_id=public_rt.id,
    destination_cidr_block="0.0.0.0/0",
    gateway_id=igw.id
)

nat_gw_route = aws.ec2.Route("nat_gw_route",
    route_table_id=private_rt.id,
    destination_cidr_block="0.0.0.0/0",
    nat_gateway_id=nat_gw.id
)

################################################################################
# Route Table Association (https://www.pulumi.com/registry/packages/aws/api-docs/ec2/routetableassociation/)
################################################################################

for i in range(len(availability_zones)):
    route_table_association = aws.ec2.RouteTableAssociation(f"public_rta_{availability_zones[i]}",
        route_table_id=public_rt.id,
        subnet_id=public_subnets[i].id
    )

for i in range(len(availability_zones)):
    route_table_association = aws.ec2.RouteTableAssociation(f"private_rta_{availability_zones[i]}",
        route_table_id=private_rt.id,
        subnet_id=private_subnets[i].id
    )

################################################################################
# EKS Cluster (https://www.pulumi.com/registry/packages/aws/api-docs/eks/cluster/)
################################################################################

eks_node_role = aws.iam.Role("eks_node_role",
    name="eks-auto-node-role",
    assume_role_policy=json.dumps({
        "Version": "2012-10-17",
        "Statement": [{
            "Action": [
                "sts:AssumeRole"
            ],
            "Effect": "Allow",
            "Principal": {
                "Service": "ec2.amazonaws.com",
            },
        }],
    })
)

eks_cluster_role = aws.iam.Role("eks_cluster_role",
    name="eks-cluster-role",
    assume_role_policy=json.dumps({
        "Version": "2012-10-17",
        "Statement": [{
            "Action": [
                "sts:AssumeRole",
                "sts:TagSession",
            ],
            "Effect": "Allow",
            "Principal": {
                "Service": "eks.amazonaws.com",
            },
        }],
    })
)

node_amazon_eks_worker_node_minimal_policy = aws.iam.RolePolicyAttachment("node_AmazonEKSWorkerNodeMinimalPolicy",
    policy_arn="arn:aws:iam::aws:policy/AmazonEKSWorkerNodeMinimalPolicy",
    role=eks_node_role.name
)

node_amazon_ec2_container_registry_pull_only = aws.iam.RolePolicyAttachment("node_AmazonEC2ContainerRegistryPullOnly",
    policy_arn="arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly",
    role=eks_node_role.name)

cluster_amazon_eks_cluster_policy = aws.iam.RolePolicyAttachment("cluster_AmazonEKSClusterPolicy",
    policy_arn="arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    role=eks_cluster_role.name
)

cluster_amazon_eks_compute_policy = aws.iam.RolePolicyAttachment("cluster_AmazonEKSComputePolicy",
    policy_arn="arn:aws:iam::aws:policy/AmazonEKSComputePolicy",
    role=eks_cluster_role.name
)

cluster_amazon_eks_block_storage_policy = aws.iam.RolePolicyAttachment("cluster_AmazonEKSBlockStoragePolicy",
    policy_arn="arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy",
    role=eks_cluster_role.name
)

cluster_amazon_eks_load_balancing_policy = aws.iam.RolePolicyAttachment("cluster_AmazonEKSLoadBalancingPolicy",
    policy_arn="arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy",
    role=eks_cluster_role.name
)

cluster_amazon_eks_networking_policy = aws.iam.RolePolicyAttachment("cluster_AmazonEKSNetworkingPolicy",
    policy_arn="arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy",
    role=eks_cluster_role.name
)

eks_cluster = aws.eks.Cluster("eks_cluster",
    name=f"eks-{pulumi.get_project()}",
    access_config={
        "authentication_mode": "API",
        "bootstrap_cluster_creator_admin_permissions": True,
    },
    role_arn=eks_cluster_role.arn,
    version="1.35",
    bootstrap_self_managed_addons=False,
    compute_config={
        "enabled": True,
        "node_pools": [
            "general-purpose",
            "system",
        ],
        "node_role_arn": eks_node_role.arn,
    },
    kubernetes_network_config={
        "elastic_load_balancing": {
            "enabled": True,
        },
    },
    storage_config={
        "block_storage": {
            "enabled": True,
        },
    },
    vpc_config={
        "endpoint_private_access": True,
        "endpoint_public_access": True,
        "subnet_ids": private_subnets,
    },
    opts = pulumi.ResourceOptions(depends_on=[
            cluster_amazon_eks_cluster_policy,
            cluster_amazon_eks_compute_policy,
            cluster_amazon_eks_block_storage_policy,
            cluster_amazon_eks_load_balancing_policy,
            cluster_amazon_eks_networking_policy,
        ]),
    tags={
        "Name": f"eks-{pulumi.get_project()}",
        "environment": pulumi.get_stack(),
    }
)