Kurtosis package indexer
========================

Kurtosis package indexer is a backend services searching for Kurtosis packages in GitHub and storing them in memory.
Right now it is consumed by Kurtosis Frontend to power Kurtosis Packages Catalog.

Implementation details
----------------------

The service simply runs a job periodically to search for all Kurtosis Packages currently existing on GitHub.
- The background job runs every two hours. Results are stored in memory for now. I.e. restarting the service will re-run the job
- It searches for `kurtosis.yml` files on GitHub. It then checks the `kurtosis.yml` file can be parsed, and there is a valid `main.star` file next to it. Any folder not matching those criteria will be discarded

GitHub authentication
---------------------

The searches run on GitHub need to be authenticated. There are two ways Kurtosis Package Indexer will authenticate itself
on GitHub.
Right now, the indexer first tries reading the `GITHUB_USER_TOKEN` environment variable and if it's empty, it falls back
to the S3 bucket option.

### Using a Github token via an environment variable
This is the simplest. The indexer expects a valid GitHub token stored inside the environment variable `GITHUB_USER_TOKEN`.

### Using a file stored inside an S3 bucket
The indexer can also get the GitHub token from a file stored inside an S3 bucket.
The file storing the GitHub token should be named `github-user-token.txt` and it should contain only the GitHub token 
on as plain text.

To access this file, the indexer will require the following environment variables to be set:
- `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` for AWS authentication
- `AWS_BUCKET_REGION` and `AWS_BUCKET_NAME` to identify the AWS S3 bucket. The user linked to the key above needs to  
be able to do `GetObject` on this bucket
- `AWS_BUCKET_FOLDER` (optional) in case the file `github-user-token.txt` is located inside a folder in this S3 bucket

Metrics authentication
----------------------

The indexer consume some Kurtosis public metrics, just package run counts for now, in order to provide this information 
to indexer clients like the package catalog.

[Snowflake][snowflake] is the Kurtosis metrics storage at the moment, and the indexer is using the [Go Snowflake client]
[gosnowflake] to execute queries on it.

It's necessary to validate a user before executing any query on this storage, we are created a new service account
and a new role for this purpose, you can access into the Kurtosis Snowflake account to get this information.

The indexer will require the following environment variables to be set:
- `KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER` for identify the Kurtosis SF account [using this format][snowflake-account-format].
- `KURTOSIS_SNOWFLAKE_DB` to specify the metrics db name
- `KURTOSIS_SNOWFLAKE_USER` the Kurtosis backend service account user
- `KURTOSIS_SNOWFLAKE_PASSWORD` the Kurtosis backend service account password
- `KURTOSIS_SNOWFLAKE_ROLE` the specific role to get access to the public metrics
- `KURTOSIS_SNOWFLAKE_WAREHOUSE` the metrics warehouse name

Data persistence
----------------

The Kurtosis packages information are stored by default in-memory. Everytime the indexer is restarted, it re-runs the
GitHub searches to fetch the latest information about the packages on GitHub.

~~There's also the option of persisting the data to a [bolt](https://github.com/etcd-io/bbolt) key value store, so that 
services can be restarted keeping the data intact. To use it, the environment variable `BOLT_DATABASE_FILE_PATH` can 
be set to point to a file on disk that bolt will use to store the data. If the indexer is being run in a container, a 
persistent volume should be used to fully benefit from this feature.~~

~~Ultimately, to make the indexer fully stateless, data can also be stored in an external 
[ETCD](https://etcd.io/) key value store. Once the ETCD cluster is up and running, the indexer can be started with the
environment variable `ETCD_DATABASE_URLS` set to the list of ETCD nodes URLs separated by a comma: 
`http://etcd.node.1:2379,http://etcd.node.2:2379,http://etcd.node.3:2379`.~~

The bolt db and the etcd db implementations were deprecated because these were not used in production so, we decided to
deprecate them in order to simplify code maintenance.

[snowflake]: https://www.snowflake.com/en/
[gosnowflake]: https://github.com/snowflakedb/gosnowflake
[snowflake-account-format]: https://docs.snowflake.com/en/user-guide/admin-account-identifier#format-1-preferred-account-name-in-your-organization