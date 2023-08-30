KURTOSIS_PACKAGE_INDEXER_IMAGE = "kurtosistech/kurtosis-package-indexer"
KURTOSIS_PACKAGE_INDEXER_DEFAULT_VERSION = "0.0.6"
KURTOSIS_PACKAGE_INDEXER_PORT_NUM = 9770

AWS_S3_BUCKET_SUBFOLDER = "kurtosis-package-indexer"


def run(
    plan,
    github_user_token="",                        # type: string
    aws_access_key_id="",                        # type: string
    aws_secret_access_key="",                    # type: string
    aws_bucket_region="",                        # type: string
    aws_bucket_name="",                          # type: string
    aws_bucket_user_folder="",                   # type: string
    kurtosis_package_indexer_custom_version="",  # type: string
):
    indexer_env_vars = {}
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
    plan.add_service(
        name="kurtosis-package-indexer",
        config=ServiceConfig(
            image=image_name_and_version,
            ports={
                "http": PortSpec(KURTOSIS_PACKAGE_INDEXER_PORT_NUM),
            },
            env_vars=indexer_env_vars,
            ready_conditions=ReadyCondition(
                recipe=PostHttpRequestRecipe(
                    port_id="http",
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
