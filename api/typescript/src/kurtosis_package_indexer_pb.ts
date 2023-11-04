// @generated by protoc-gen-es v1.4.0 with parameter "target=ts"
// @generated from file kurtosis_package_indexer.proto (package kurtosis_package_indexer, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64, Timestamp } from "@bufbuild/protobuf";

/**
 * @generated from enum kurtosis_package_indexer.ArgumentValueType
 */
export enum ArgumentValueType {
  /**
   * @generated from enum value: BOOL = 0;
   */
  BOOL = 0,

  /**
   * @generated from enum value: STRING = 1;
   */
  STRING = 1,

  /**
   * @generated from enum value: INTEGER = 2;
   */
  INTEGER = 2,

  /**
   * @generated from enum value: DICT = 4;
   */
  DICT = 4,

  /**
   * @generated from enum value: JSON = 5;
   */
  JSON = 5,

  /**
   * @generated from enum value: LIST = 6;
   */
  LIST = 6,
}
// Retrieve enum metadata with: proto3.getEnumType(ArgumentValueType)
proto3.util.setEnumType(ArgumentValueType, "kurtosis_package_indexer.ArgumentValueType", [
  { no: 0, name: "BOOL" },
  { no: 1, name: "STRING" },
  { no: 2, name: "INTEGER" },
  { no: 4, name: "DICT" },
  { no: 5, name: "JSON" },
  { no: 6, name: "LIST" },
]);

/**
 * @generated from message kurtosis_package_indexer.ReadPackageRequest
 */
export class ReadPackageRequest extends Message<ReadPackageRequest> {
  /**
   * @generated from field: kurtosis_package_indexer.PackageRepository repository_metadata = 1;
   */
  repositoryMetadata?: PackageRepository;

  constructor(data?: PartialMessage<ReadPackageRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.ReadPackageRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "repository_metadata", kind: "message", T: PackageRepository },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ReadPackageRequest {
    return new ReadPackageRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ReadPackageRequest {
    return new ReadPackageRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ReadPackageRequest {
    return new ReadPackageRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ReadPackageRequest | PlainMessage<ReadPackageRequest> | undefined, b: ReadPackageRequest | PlainMessage<ReadPackageRequest> | undefined): boolean {
    return proto3.util.equals(ReadPackageRequest, a, b);
  }
}

/**
 * @generated from message kurtosis_package_indexer.ReadPackageResponse
 */
export class ReadPackageResponse extends Message<ReadPackageResponse> {
  /**
   * @generated from field: optional kurtosis_package_indexer.KurtosisPackage package = 1;
   */
  package?: KurtosisPackage;

  constructor(data?: PartialMessage<ReadPackageResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.ReadPackageResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "package", kind: "message", T: KurtosisPackage, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ReadPackageResponse {
    return new ReadPackageResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ReadPackageResponse {
    return new ReadPackageResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ReadPackageResponse {
    return new ReadPackageResponse().fromJsonString(jsonString, options);
  }

  static equals(a: ReadPackageResponse | PlainMessage<ReadPackageResponse> | undefined, b: ReadPackageResponse | PlainMessage<ReadPackageResponse> | undefined): boolean {
    return proto3.util.equals(ReadPackageResponse, a, b);
  }
}

/**
 * @generated from message kurtosis_package_indexer.GetPackagesResponse
 */
export class GetPackagesResponse extends Message<GetPackagesResponse> {
  /**
   * @generated from field: repeated kurtosis_package_indexer.KurtosisPackage packages = 1;
   */
  packages: KurtosisPackage[] = [];

  constructor(data?: PartialMessage<GetPackagesResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.GetPackagesResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "packages", kind: "message", T: KurtosisPackage, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetPackagesResponse {
    return new GetPackagesResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetPackagesResponse {
    return new GetPackagesResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetPackagesResponse {
    return new GetPackagesResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetPackagesResponse | PlainMessage<GetPackagesResponse> | undefined, b: GetPackagesResponse | PlainMessage<GetPackagesResponse> | undefined): boolean {
    return proto3.util.equals(GetPackagesResponse, a, b);
  }
}

/**
 * @generated from message kurtosis_package_indexer.KurtosisPackage
 */
export class KurtosisPackage extends Message<KurtosisPackage> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: repeated kurtosis_package_indexer.PackageArg args = 2;
   */
  args: PackageArg[] = [];

  /**
   * @generated from field: uint64 stars = 3;
   */
  stars = protoInt64.zero;

  /**
   * @generated from field: string description = 4;
   */
  description = "";

  /**
   * deprecated: use a combination of repository_url and root_path instead
   *
   * @generated from field: optional string url = 5;
   */
  url?: string;

  /**
   * @generated from field: string entrypoint_description = 6;
   */
  entrypointDescription = "";

  /**
   * @generated from field: string returns_description = 7;
   */
  returnsDescription = "";

  /**
   * @generated from field: kurtosis_package_indexer.PackageRepository repository_metadata = 8;
   */
  repositoryMetadata?: PackageRepository;

  /**
   * @generated from field: string parsing_result = 9;
   */
  parsingResult = "";

  /**
   * @generated from field: google.protobuf.Timestamp parsing_time = 10;
   */
  parsingTime?: Timestamp;

  /**
   * @generated from field: string version = 11;
   */
  version = "";

  constructor(data?: PartialMessage<KurtosisPackage>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.KurtosisPackage";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "args", kind: "message", T: PackageArg, repeated: true },
    { no: 3, name: "stars", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 4, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 6, name: "entrypoint_description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "returns_description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 8, name: "repository_metadata", kind: "message", T: PackageRepository },
    { no: 9, name: "parsing_result", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 10, name: "parsing_time", kind: "message", T: Timestamp },
    { no: 11, name: "version", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): KurtosisPackage {
    return new KurtosisPackage().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): KurtosisPackage {
    return new KurtosisPackage().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): KurtosisPackage {
    return new KurtosisPackage().fromJsonString(jsonString, options);
  }

  static equals(a: KurtosisPackage | PlainMessage<KurtosisPackage> | undefined, b: KurtosisPackage | PlainMessage<KurtosisPackage> | undefined): boolean {
    return proto3.util.equals(KurtosisPackage, a, b);
  }
}

/**
 * @generated from message kurtosis_package_indexer.PackageArg
 */
export class PackageArg extends Message<PackageArg> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: bool is_required = 2;
   */
  isRequired = false;

  /**
   * @generated from field: string description = 4;
   */
  description = "";

  /**
   * @generated from field: kurtosis_package_indexer.PackageArgumentType typeV2 = 5;
   */
  typeV2?: PackageArgumentType;

  /**
   * @generated from field: optional string defaultValue = 6;
   */
  defaultValue?: string;

  constructor(data?: PartialMessage<PackageArg>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.PackageArg";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "is_required", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 4, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "typeV2", kind: "message", T: PackageArgumentType },
    { no: 6, name: "defaultValue", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PackageArg {
    return new PackageArg().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PackageArg {
    return new PackageArg().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PackageArg {
    return new PackageArg().fromJsonString(jsonString, options);
  }

  static equals(a: PackageArg | PlainMessage<PackageArg> | undefined, b: PackageArg | PlainMessage<PackageArg> | undefined): boolean {
    return proto3.util.equals(PackageArg, a, b);
  }
}

/**
 * @generated from message kurtosis_package_indexer.PackageArgumentType
 */
export class PackageArgumentType extends Message<PackageArgumentType> {
  /**
   * @generated from field: kurtosis_package_indexer.ArgumentValueType top_level_type = 1;
   */
  topLevelType = ArgumentValueType.BOOL;

  /**
   * @generated from field: optional kurtosis_package_indexer.ArgumentValueType inner_type_1 = 2;
   */
  innerType1?: ArgumentValueType;

  /**
   * @generated from field: optional kurtosis_package_indexer.ArgumentValueType inner_type_2 = 3;
   */
  innerType2?: ArgumentValueType;

  constructor(data?: PartialMessage<PackageArgumentType>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.PackageArgumentType";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "top_level_type", kind: "enum", T: proto3.getEnumType(ArgumentValueType) },
    { no: 2, name: "inner_type_1", kind: "enum", T: proto3.getEnumType(ArgumentValueType), opt: true },
    { no: 3, name: "inner_type_2", kind: "enum", T: proto3.getEnumType(ArgumentValueType), opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PackageArgumentType {
    return new PackageArgumentType().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PackageArgumentType {
    return new PackageArgumentType().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PackageArgumentType {
    return new PackageArgumentType().fromJsonString(jsonString, options);
  }

  static equals(a: PackageArgumentType | PlainMessage<PackageArgumentType> | undefined, b: PackageArgumentType | PlainMessage<PackageArgumentType> | undefined): boolean {
    return proto3.util.equals(PackageArgumentType, a, b);
  }
}

/**
 * @generated from message kurtosis_package_indexer.PackageRepository
 */
export class PackageRepository extends Message<PackageRepository> {
  /**
   * @generated from field: string base_url = 1;
   */
  baseUrl = "";

  /**
   * @generated from field: string owner = 2;
   */
  owner = "";

  /**
   * @generated from field: string name = 3;
   */
  name = "";

  /**
   * @generated from field: string root_path = 4;
   */
  rootPath = "";

  constructor(data?: PartialMessage<PackageRepository>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "kurtosis_package_indexer.PackageRepository";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "base_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "owner", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "root_path", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PackageRepository {
    return new PackageRepository().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PackageRepository {
    return new PackageRepository().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PackageRepository {
    return new PackageRepository().fromJsonString(jsonString, options);
  }

  static equals(a: PackageRepository | PlainMessage<PackageRepository> | undefined, b: PackageRepository | PlainMessage<PackageRepository> | undefined): boolean {
    return proto3.util.equals(PackageRepository, a, b);
  }
}

