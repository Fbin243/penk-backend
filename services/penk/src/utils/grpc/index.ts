import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";

import type { ProtoGrpcType as CoreGrpcType } from "./proto/core_service";

const coreService = grpc.loadPackageDefinition(
  protoLoader.loadSync("../../proto/core/core_service.proto"),
) as unknown as CoreGrpcType;

export const coreClient = new coreService.core.Core(
  `localhost:${process.env.CORE_GRPC_PORT || "50051"}`,
  grpc.credentials.createInsecure(),
);
