KURTOSIS_PACKAGE_INDEXER_IMAGE = "kurtosistech/kurtosis-package-indexer"
KURTOSIS_PACKAGE_INDEXER_DEFAULT_VERSION = "0.0.7"

KURTOSIS_PACKAGE_INDEXER_PORT_ID = "http"
KURTOSIS_PACKAGE_INDEXER_PORT_NUM = 9770

AWS_S3_BUCKET_SUBFOLDER = "kurtosis-package-indexer"


def run(
    plan,
    github_user_token="",
    aws_access_key_id="",
    aws_secret_access_key="",
    aws_bucket_region="",
    aws_bucket_name="",
    aws_bucket_user_folder="",
    kurtosis_package_indexer_custom_version="",
):
    """Runs a Kurtosis package indexer service, listening on port 9770

    Args:
        github_user_token: The Github user token to use to authenticate to Github
        aws_access_key_id: The AWS access key ID to authenticate to AWS
        aws_secret_access_key: The AWS secret key to authenticate to AWS
        aws_bucket_region: The region of the AWS bucket to read from
        aws_bucket_name: The name of the bucket to read from
        aws_bucket_user_folder: The folder inside the AWS bucket
        kurtosis_package_indexer_custom_version: A custom version for the container image
    Returns:
        The service object containing useful information on the running service. Typically:
        ```
        {
            "hostname": "kurtosis-package-indexer",
            "ip_address": "172.16.0.4",
            "name": "kurtosis-package-indexer",
            "port_number": 9770
        }
        ```
    """
    indexer_env_vars = {
        "BOLT_DATABASE_FILE_PATH": "/data/bolt.db"
    }
    if len(github_user_token) > 0:
        indexer_env_vars["GITHUB_USER_TOKEN"] = github_user_token
    else:
        aws_env = get_aws_env(aws_access_key_id, aws_secret_access_key, aws_bucket_region, aws_bucket_name, aws_bucket_user_folder)
        indexer_env_vars["AWS_ACCESS_KEY_ID"] = aws_env.access_key_id
        indexer_env_vars["AWS_SECRET_ACCESS_KEY"] = aws_env.secret_access_key
        indexer_env_vars["AWS_BUCKET_REGION"] = aws_env.bucket_region
        indexer_env_vars["AWS_BUCKET_NAME"] = aws_env.bucket_name
        indexer_env_vars["AWS_BUCKET_FOLDER"] = "{}/{}".format(aws_env.bucket_user_folder, AWS_S3_BUCKET_SUBFOLDER)

    image_name_and_version = "{}:{}".format(
        KURTOSIS_PACKAGE_INDEXER_IMAGE,
        kurtosis_package_indexer_custom_version if kurtosis_package_indexer_custom_version != "" else KURTOSIS_PACKAGE_INDEXER_DEFAULT_VERSION
    )
    indexer_service = plan.add_service(
        name="kurtosis-package-indexer",
        config=ServiceConfig(
            image=image_name_and_version,
            ports={
                KURTOSIS_PACKAGE_INDEXER_PORT_ID: PortSpec(KURTOSIS_PACKAGE_INDEXER_PORT_NUM),
            },
            files={
                "/data/": Directory(
                    persistent_key="bold_db_data"
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



def get_aws_env(aws_access_key_id, aws_secret_access_key, aws_bucket_region, aws_bucket_name, aws_bucket_user_folder):
    if len(aws_access_key_id) == 0 and len(aws_secret_access_key) == 0 and len(aws_bucket_region) and len(aws_bucket_name) == 0:
        return struct(
            access_key_id=aws_access_key_id,
            secret_access_key=aws_secret_access_key,
            bucket_region=aws_bucket_region,
            bucket_name=aws_bucket_name,
            bucket_user_folder=aws_bucket_user_folder,
        )
    # the AWS values should be provided as env variables to the package. Otherwise this package cannot run
    return struct(
        access_key_id=kurtosis.aws_access_key_id,
        secret_access_key=kurtosis.aws_secret_access_key,
        bucket_region=kurtosis.aws_bucket_region,
        bucket_name=kurtosis.aws_bucket_name,
        bucket_user_folder=kurtosis.aws_bucket_user_folder,
    )
