Kurtosis package for Kurtosis Package Indexer
=============================================

This Kurtosis package spins up a Kurtosis Package Indexer inside a Kurtosis enclave, listening on port 9770.

Running the package
-------------------

Kurtosis Package Indexer required a valid Github token to search Kurtosis packages on Github.
The token can be passed either as a direct environment variable, or the indexer can also fetch
it from the content of a file inside an S3 bukcet. See the indexer [README](../README.md) for more info 

If the latter is used, make sure a file named `github-user-token.txt` is available in the S3 bucket
under the path `${aws_bucket_user_folder}/kurtosis-package-indexer/`

The following arguments that can be passed to the package:
```json
{
    // If the service should authenticate to Github with a token, it can be passed here and the
    // aws_* other args can be ignored
    "github_user_token": "",         // optional

    // If it is expected that the service will get the Github user token from an S3 bucket, those
    // args should be filled. `aws_bucket_user_folder` can remain empty is the file containing the
    // token is at the root of the bucket
    "aws_access_key_id": "",        // optional
    "aws_secret_access_key": "",    // optional
    "aws_bucket_region": "",        // optional
    "aws_bucket_name": "",          // optional
    "aws_bucket_user_folder": ""    // optional
}
```

Note that when running this package on Kurtosis cloud, the package will naturally use the AWS environment variable
automatically provided to the package to fetch the Github token inside AWS S3.
