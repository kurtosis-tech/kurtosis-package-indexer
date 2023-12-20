etcd_module = import_module("github.com/kurtosis-tech/etcd-package/main.star")

KURTOSIS_PACKAGE_INDEXER_IMAGE = "kurtosistech/kurtosis-package-indexer"

KURTOSIS_PACKAGE_INDEXER_DATA_DIR = "/data"
KURTOSIS_PACKAGE_INDEXER_PORT_ID = "http"
KURTOSIS_PACKAGE_INDEXER_PORT_NUM = 9770

AWS_S3_BUCKET_SUBFOLDER = "kurtosis-package-indexer"

# Kurtosis Snowflake Account Keys
KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER_KEY = "kurtosis_snowflake_account_identifier"
KURTOSIS_SNOWFLAKE_DB_KEY = "kurtosis_snowflake_db"
KURTOSIS_SNOWFLAKE_PASSWORD_KEY = "kurtosis_snowflake_password"
KURTOSIS_SNOWFLAKE_ROLE_KEY = "kurtosis_snowflake_role"
KURTOSIS_SNOWFLAKE_USER_KEY = "kurtosis_snowflake_user"
KURTOSIS_SNOWFLAKE_WAREHOUSE_KEY = "kurtosis_snowflake_warehouse"


def run(
    plan,
    github_user_token="",
    kurtosis_package_indexer_version="0.0.7",
    snowflake_env={},
    aws_env={},
):
    """Runs a Kurtosis package indexer service, listening on port 9770

    Args:
        github_user_token (string): The GitHub user token to use to authenticate to GitHub
        kurtosis_package_indexer_version (string): The version of the container image to use
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

    indexer_env_vars = get_snowflake_env(snowflake_env)
    if len(github_user_token) > 0:
        indexer_env_vars["GITHUB_USER_TOKEN"] = github_user_token
    else:
        aws_env = get_aws_env(aws_env)
        indexer_env_vars["AWS_ACCESS_KEY_ID"] = aws_env.access_key_id
        indexer_env_vars["AWS_SECRET_ACCESS_KEY"] = aws_env.secret_access_key
        indexer_env_vars["AWS_BUCKET_REGION"] = aws_env.bucket_region
        indexer_env_vars["AWS_BUCKET_NAME"] = aws_env.bucket_name
        indexer_env_vars["AWS_BUCKET_FOLDER"] = "{}/{}".format(aws_env.bucket_user_folder, AWS_S3_BUCKET_SUBFOLDER)

    image_name_and_version = "{}:{}".format(KURTOSIS_PACKAGE_INDEXER_IMAGE, kurtosis_package_indexer_version)
    indexer_service = plan.add_service(
        name="kurtosis-package-indexer",
        config=ServiceConfig(
            image=image_name_and_version,
            ports={
                KURTOSIS_PACKAGE_INDEXER_PORT_ID: PortSpec(KURTOSIS_PACKAGE_INDEXER_PORT_NUM),
            },
            files={
                KURTOSIS_PACKAGE_INDEXER_DATA_DIR: Directory(
                    persistent_key="kurtosis_indexer_data_dir"
                ),
            },
            env_vars=indexer_env_vars,
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
        port_number=indexer_service.ports[KURTOSIS_PACKAGE_INDEXER_PORT_ID].number
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
    env_vars={}
    # the Snowflake values should be provided as env variables to the package, or  the "CI" env var has to be true to use the `doNothingReporter`
    sf_account_identifier = sf_env[KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER_KEY] if KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER_KEY in sf_env else ""
    sf_db = sf_env[KURTOSIS_SNOWFLAKE_DB_KEY] if KURTOSIS_SNOWFLAKE_DB_KEY in sf_env else ""
    sf_password = sf_env[KURTOSIS_SNOWFLAKE_PASSWORD_KEY] if KURTOSIS_SNOWFLAKE_PASSWORD_KEY in sf_env else ""
    sf_role = sf_env[KURTOSIS_SNOWFLAKE_ROLE_KEY] if KURTOSIS_SNOWFLAKE_ROLE_KEY in sf_env else ""
    sf_user = sf_env[KURTOSIS_SNOWFLAKE_USER_KEY] if KURTOSIS_SNOWFLAKE_USER_KEY in sf_env else ""
    sf_warehouse = sf_env[KURTOSIS_SNOWFLAKE_WAREHOUSE_KEY] if KURTOSIS_SNOWFLAKE_WAREHOUSE_KEY in sf_env else ""

    env_vars["KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER"]=sf_account_identifier
    env_vars["KURTOSIS_SNOWFLAKE_DB"]=sf_db
    env_vars["KURTOSIS_SNOWFLAKE_PASSWORD"]=sf_password
    env_vars["KURTOSIS_SNOWFLAKE_ROLE"]=sf_role
    env_vars["KURTOSIS_SNOWFLAKE_USER"]=sf_user
    env_vars["KURTOSIS_SNOWFLAKE_WAREHOUSE"]=sf_warehouse

    return env_vars
