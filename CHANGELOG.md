# Changelog

## [0.0.30](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.29...0.0.30) (2024-03-01)


### Features

* add the package `locator_root` in the responses ([#116](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/116)) ([a1ec7c6](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/a1ec7c659481dc75f6e2bf2618b636b5bb4a4922)), closes [#115](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/115)

## [0.0.29](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.28...0.0.29) (2024-02-26)


### Bug Fixes

* The docker compose package arguments should all be optional ([#113](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/113)) ([d679d9d](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d679d9da55ff6a489f0121ec587caaa4f242620e))

## [0.0.28](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.27...0.0.28) (2024-02-23)


### Features

* added a method that works on raw file content ([#109](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/109)) ([6d21bec](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/6d21bec492d50bf4398f4f75453e8f5d0ac12a19))
* Extract package arguments from the docker compose env file ([#112](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/112)) ([3c2b5e3](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/3c2b5e3c8c8cc1cdbc67f043544e49c5ec4a47ce))

## [0.0.27](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.26...0.0.27) (2024-02-02)


### Features

* configuring logger log level at runtime ([#107](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/107)) ([bd87d35](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/bd87d35a7262163f3c745a2eece308913219019c))

## [0.0.26](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.25...0.0.26) (2024-01-26)


### Bug Fixes

* add compose to read pkg err ([#104](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/104)) ([7fa417a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/7fa417afd1642d16a595fecf4523c8b6e6eec191))
* Docker compose filepath creation in the extractDockerComposePackageContent function ([#106](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/106)) ([1414a0c](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/1414a0c98f1691c37104b8abae5b5fd48bf60137))

## [0.0.25](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.24...0.0.25) (2024-01-19)


### Bug Fixes

* adding a new line in the Readme file to force release please trigger ([#101](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/101)) ([40a07db](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/40a07db8754d2b6dc51ed394b12bb7897290f2ea))

## [0.0.24](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.23...0.0.24) (2024-01-18)


### Features

* support reading from docker compose based packages ([#60](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/60)) ([57a7564](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/57a75645cf0a38e37a7d3910bf0a724648f37c8a))

## [0.0.23](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.22...0.0.23) (2024-01-09)


### Bug Fixes

* crawler accidentally stopped ([#96](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/96)) ([71a7d0e](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/71a7d0e06d04e7bc83600f53e85d70e0cd1a4131))

## [0.0.22](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.21...0.0.22) (2024-01-09)


### Bug Fixes

* fails with dict argument, in the doc string, has not been parameterized ([#94](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/94)) ([87666e4](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/87666e4c03b8fe00c72cfdf90d7805ec78851052))

## [0.0.21](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.20...0.0.21) (2023-12-20)


### Features

* running a secondary GitHub crawler for updating the repository stars and the last commit time with a longer frequency ([#87](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/87)) ([caab53d](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/caab53d372e82267242b5bf93e5faa53b5b8ce63))

## [0.0.20](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.19...0.0.20) (2023-12-18)


### Features

* reading packages catalog from yaml file ([#77](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/77)) ([d3ef0bd](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d3ef0bd31cabe13fdbb3782a5b4c0548998780c8))

## [0.0.19](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.18...0.0.19) (2023-12-15)


### Features

* remove bolt and etcd storage implementations ([#83](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/83)) ([b94c6cc](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/b94c6cc15be4c275897b729017ff9b15c6f9e8df))

## [0.0.18](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.17...0.0.18) (2023-12-14)


### Features

* publish the autogenerated TS bindings ([#81](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/81)) ([fb3332b](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/fb3332bfac0ef23c564d20d24cd164899c8cad67))

## [0.0.17](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.16...0.0.17) (2023-12-12)


### Features

* add packages run count metric ([#74](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/74)) ([09b8ded](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/09b8ded18455ed044e1402137d9c65e7812f0b6b))

## [0.0.16](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.15...0.0.16) (2023-12-04)


### Bug Fixes

* returning icon URL for Kurtosis packages in subfolders ([#72](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/72)) ([5bddfad](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/5bddfad6dc66ce6ee6f9bbee554c20474733a2a3))

## [0.0.15](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.14...0.0.15) (2023-12-04)


### Features

* kurtosis package icon URL added in the returned package info ([#70](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/70)) ([5015cdb](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/5015cdbc733fea9262a8f1c27804435c6fde1a1d))

## [0.0.14](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.13...0.0.14) (2023-12-01)


### Features

* return the latest commit data to a package ([#68](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/68)) ([72be6c6](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/72be6c69e2c1012e7a69567ac1a8e8507e82d66d)), closes [#65](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/65)

## [0.0.13](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.12...0.0.13) (2023-11-30)


### Bug Fixes

* revet change on repo filepaths generation ([#66](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/66)) ([1b5bb45](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/1b5bb45513bbf94c7a247278c74526691c6a5429))

## [0.0.12](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.11...0.0.12) (2023-11-30)


### Bug Fixes

* returns the package starts ([#64](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/64)) ([c5e077a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/c5e077aabcc4795c0ba5be51a2a036f4f82027c6))

## [0.0.11](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.10...0.0.11) (2023-11-07)


### Features

* generate typescript bindings ([#46](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/46)) ([bc881c5](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/bc881c5be5d424d1a6a875350e43804a33f1c0d7))
* return default values ([#54](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/54)) ([21739bc](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/21739bc35b3c3dd4aa44976c7793e6e4914887b6))


### Bug Fixes

* improve read pkg err handling ([#56](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/56)) ([4c090ff](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/4c090ffd9235eed18f0dee20b9941a231367adbf))
* use forked kurtosis docstring pkg ([#53](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/53)) ([7a175e6](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/7a175e6f62d4fa249b03cccb199465691a7e3b4e))

## [0.0.10](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.9...0.0.10) (2023-10-19)


### Features

* get package on demand ([#42](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/42)) ([8094c1d](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/8094c1d02a6b7dc2cc5b82658b7c540af9363e67))

## [0.0.9](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.8...0.0.9) (2023-09-21)


### Features

* Add ETCD as an option for the store ([#32](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/32)) ([557fb51](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/557fb51f9e97707471c284d07581a9db0f4866a6))
* Filter out packages with invalid name in `kurtosis.yml` ([#37](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/37)) ([df53368](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/df53368d864b5811506f67634e0a5f7716e40633))
* Service is automatically deployed to AWS on release ([#38](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/38)) ([93bf876](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/93bf87612a68d21d4dfb799f7e9f6b5edfcf1e22))

## [0.0.8](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.7...0.0.8) (2023-09-15)


### Features

* Add Kurtosis package to run the indexer in a Kurtosis enclave ([#19](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/19)) ([c343f0d](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/c343f0d70c92c65eae1c5cac39bea3f5e03ca293))
* Parse docstring of the run function in main.star ([#28](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/28)) ([09b6692](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/09b669296b3a8248c14e04e4a5fc827f2d007654))
* Persistent KV store using BBolt ([#31](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/31)) ([5748f25](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/5748f253910369c5d2c8e007d2c07af8fa648355))
* Support `list` type in docstrings ([#30](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/30)) ([3b778a4](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/3b778a4973efe843d30c7145b21a8e0f64246688))

## [0.0.7](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.6...0.0.7) (2023-08-31)


### Bug Fixes

* Reindex endpoint was failing b/c the context used was canceled after the endpoint returned ([#24](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/24)) ([f6b1897](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/f6b1897252c9bbd32237a589f56dfcd7578ffb8b))

## [0.0.6](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.5...0.0.6) (2023-08-31)


### Features

* Add possibility to retrieve Github Token from the user S3 bucket ([#18](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/18)) ([b1924f4](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/b1924f453f68b7dc056af52fa9c3a48e01c0e7af))

## [0.0.5](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.4...0.0.5) (2023-08-30)


### Features

* Add endpoint to manually trigger a reindex ([#20](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/20)) ([d6c777a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d6c777a573329ee16389fb007ed1880b7e6ec1ec))
* Add package URL in returned Kurtosis package objects ([#17](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/17)) ([d1f6e6b](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d1f6e6b294e93b3792bd8a69df4c28df9dc93c14))
* Parse package description from kurtosis.yml ([#21](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/21)) ([65d6e81](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/65d6e811f4b7f62c34ef8dcf449992638b563bd5))

## [0.0.4](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.3...0.0.4) (2023-08-29)


### Bug Fixes

* Package with a `kurtosis.yml` file with an empty name are not considered invalid ([#15](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/15)) ([6b989f4](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/6b989f4c6832012603ae4686e92f8c1524b840f1))

## [0.0.3](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.2...0.0.3) (2023-08-29)


### Features

* Improve Kurtosis packages search ([#10](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/10)) ([a1818ca](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/a1818ca1cbd990e845ae7bcd79af0a3dc58f83be))

## [0.0.2](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.1...0.0.2) (2023-08-28)


### Features

* Add connect-go for the FE to be able to talk to the indexer ([#5](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/5)) ([d796aa3](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d796aa39c4b1d1d3c00c5ed2ab6defd4338245d5))
* Add GetPackages endpoint ([#1](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/1)) ([9b9f63a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/9b9f63a9ff7bbd99289650e2a962ad9163fdf102))
* Add Github auth based on user token ([#7](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/7)) ([b2d4c4a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/b2d4c4a762a35e1faf46acc440677ffcdbcb7882))
* Crawl Github to fetch existing Kurtosis packages ([#4](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/4)) ([bc8d573](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/bc8d573036cfbb47f86120132a49c2980431f921))
* Return Github repo stars as well ([#8](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/8)) ([c9fbe5d](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/c9fbe5d22f19457b06c5971c66441416f123467b))
* The crawler now parses the main.star file ([#6](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/6)) ([3f09acb](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/3f09acbf1cdab73dcbe7ba8ab987d327abdc68bd))
