etcd_module = import_module("github.com/kurtosis-tech/etcd-package/main.star")

KURTOSIS_PACKAGE_INDEXER_IMAGE = "kurtosistech/kurtosis-package-indexer"
KURTOSIS_PACKAGE_INDEXER_DATA_DIR = "/data"
KURTOSIS_PACKAGE_INDEXER_PORT_ID = "http"
KURTOSIS_PACKAGE_INDEXER_PORT_NUM = 9770
AWS_S3_BUCKET_SUBFOLDER = "kurtosis-package-indexer"

def run(
    plan,
    is_running_in_prod="false",
    github_user_token="",
    kurtosis_package_indexer_version="",
    snowflake_env={},
    aws_env={},
):
    """Runs a Kurtosis package indexer service, listening on port 9770

    Args:
        is_running_in_prod(string): Set to false if devving locally or in CI. (this will ignore any snowflake settings)
        github_user_token (string): The GitHub user token to use to authenticate to GitHub. If empty, uses aws_env values set to retrieve a token from S3.
        kurtosis_package_indexer_version (string): The version of the container image to use. If not specified, the local code will be used to build an image.
        snowflake_env (dict[string, string]): The Snowflake information required to connect to the Snowflake account
            to get some package metrics from this storage.
            The dictionary should contain the following fields:
            ```
            {
                "kurtosis_snowflake_account_identifier": "<KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER>",
                "kurtosis_snowflake_db": "<KURTOSIS_SNOWFLAKE_DB>",
                "kurtosis_snowflake_password": "<KURTOSIS_SNOWFLAKE_PASSWORD>",
                "kurtosis_snowflake_role": "<KURTOSIS_SNOWFLAKE_ROLE>",
                "kurtosis_snowflake_user": "<KURTOSIS_SNOWFLAKE_USER>",
                "kurtosis_snowflake_warehouse": "<KURTOSIS_SNOWFLAKE_WAREHOUSE>",
            }
            ```
        aws_env (dict[string, string]): The AWS information required to optionally pull the GitHub token
            from a file in an AWS bucket. This file should be located at 
            `<BUCKET_ROOT>/<OPTIONAL_FOLDER>/kurtosis-package-indexer/github-user-token.txt`. 
            The dictionary should contain the following fields:
            ```
            {
                "aws_access_key_id": "<AWS_KEY_ID_TO_AUTHENTICATE>",
                "aws_secret_access_key": "<AWS_SECRET_ACCESS_KEY_TO_AUTHENTICATE>",
                "aws_bucket_region": "<AWS_BUCKET_REGION>",
                "aws_bucket_name": "<AWS_BUCKET_NAME>",
                "aws_bucket_user_folder": "<OPTIONAL_FOLDER_IN_AWS_BUCKET>",
            }
            ```
    Returns:
        The service object containing useful information on the running Kurtosis Package Indexer. Typically:
        ```
        {
            "hostname": "kurtosis-package-indexer",
            "ip_address": "172.16.0.4",
            "name": "kurtosis-package-indexer",
            "port_number": 9770
        }
        ```
    """
    snowflake_env = get_snowflake_env(snowflake_env)
    aws_env = get_aws_env(aws_env) 

    indexer_service = plan.add_service(
        name="kurtosis-package-indexer",
        config=ServiceConfig(
            # use local built docker image if no container image version specified
            image=ImageBuildSpec(image_name="kurtosis-package-indexer", build_context_dir=".") if kurtosis_package_indexer_version == "" else "{}:{}".format(KURTOSIS_PACKAGE_INDEXER_IMAGE, kurtosis_package_indexer_version),
            ports={
                KURTOSIS_PACKAGE_INDEXER_PORT_ID: PortSpec(KURTOSIS_PACKAGE_INDEXER_PORT_NUM),
            },
            files={
                KURTOSIS_PACKAGE_INDEXER_DATA_DIR: Directory(
                    persistent_key="kurtosis-indexer-data-dir"
                ),
            },
            env_vars= {
                "LOGGER_LOG_LEVEL": "debug",
                "GITHUB_USER_TOKEN": github_user_token, # if empty, AWS values are used to retrieve a github token from s3
                "AWS_ACCESS_KEY_ID": aws_env.access_key_id,
                "AWS_SECRET_ACCESS_KEY": aws_env.secret_access_key,
                "AWS_BUCK={}ET_REGION": aws_env.bucket_region,
                "AWS_BUCKET_NAME": aws_env.bucket_name,
                "AWS_BUCKET_FOLDER": "{}/{}".format(aws_env.bucket_user_folder, AWS_S3_BUCKET_SUBFOLDER),
                "CI": "false" if is_running_in_prod == "true" else "true", # if true, the snowflake env values will be ignored
                "KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER": snowflake_env.account_identifier,
                "KURTOSIS_SNOWFLAKE_DB":snowflake_env.db,
                "KURTOSIS_SNOWFLAKE_PASSWORD": snowflake_env.password,
                "KURTOSIS_SNOWFLAKE_ROLE": snowflake_env.role,
                "KURTOSIS_SNOWFLAKE_USER": snowflake_env.user,
                "KURTOSIS_SNOWFLAKE_WAREHOUSE": snowflake_env.warehouse,
            },
            ready_conditions=ReadyCondition(
                recipe=PostHttpRequestRecipe(
                    port_id=KURTOSIS_PACKAGE_INDEXER_PORT_ID,
                    endpoint="/kurtosis_package_indexer.KurtosisPackageIndexer/IsAvailable",
                    content_type="application/json",
                    body="{}"
                ),
                field="code",
                assertion="==",
                target_value=200,
            )
        )
    )
    return struct(
        hostname=indexer_service.hostname,
        ip_address=indexer_service.ip_address,
        name=indexer_service.name,
    )


def get_aws_env(aws_env):
    aws_access_key_id = aws_env["access_key_id"] if "access_key_id" in aws_env else ""
    aws_secret_access_key = aws_env["secret_access_key"] if "secret_access_key" in aws_env else ""
    aws_bucket_region = aws_env["bucket_region"] if "bucket_region" in aws_env else ""
    aws_bucket_name = aws_env["bucket_name"] if "bucket_name" in aws_env else ""
    aws_bucket_user_folder = aws_env["bucket_user_folder"] if "bucket_user_folder" in aws_env else ""
    if len(aws_access_key_id) == 0 and len(aws_secret_access_key) == 0 and len(aws_bucket_region) and len(aws_bucket_name) == 0:
        # the AWS values should be provided as env variables to the package. Otherwise this package cannot run
        return struct(
            access_key_id=kurtosis.aws_access_key_id,
            secret_access_key=kurtosis.aws_secret_access_key,
            bucket_region=kurtosis.aws_bucket_region,
            bucket_name=kurtosis.aws_bucket_name,
            bucket_user_folder=kurtosis.aws_bucket_user_folder,
        )
    
    return struct(
        access_key_id=aws_access_key_id,
        secret_access_key=aws_secret_access_key,
        bucket_region=aws_bucket_region,
        bucket_name=aws_bucket_name,
        bucket_user_folder=aws_bucket_user_folder,
    )

def get_snowflake_env(sf_env):
     # the Snowflake values should be provided as env variables to the package
    sf_account_identifier = sf_env["kurtosis_snowflake_account_identifier"] if "kurtosis_snowflake_account_identifier" in sf_env else ""
    sf_db = sf_env["kurtosis_snowflake_db"] if "kurtosis_snowflake_db" in sf_env else ""
    sf_password = sf_env["kurtosis_snowflake_password"] if "kurtosis_snowflake_password" in sf_env else ""
    sf_role = sf_env["kurtosis_snowflake_user"] if "kurtosis_snowflake_user" in sf_env else ""
    sf_user = sf_env["kurtosis_snowflake_role"] if "kurtosis_snowflake_role" in sf_env else ""
    sf_warehouse = sf_env["kurtosis_snowflake_warehouse"] if "kurtosis_snowflake_warehouse" in sf_env else ""

    return struct(
        account_identifier=sf_account_identifier,
        db=sf_db,
        password=sf_password,
        role=sf_role,
        user=sf_user,
        warehouse=sf_warehouse,
    )
