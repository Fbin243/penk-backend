import * as grpc from "@grpc/grpc-js";
import { Metadata } from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";

import type { ProtoGrpcType as CoreGrpcType } from "./proto/core_service";

const coreService = grpc.loadPackageDefinition(
  protoLoader.loadSync("../../proto/core/core_service.proto", {
    includeDirs: ["../../proto"],
  }),
) as unknown as CoreGrpcType;

export const coreClient = new coreService.core.Core(
  `localhost:${process.env.CORE_GRPC_PORT || "50051"}`,
  grpc.credentials.createInsecure(),
);

export const createMetadata = (firebaseUID: string): Metadata => {
  const metadata = new Metadata();
  metadata.add("service-name", "penk");
  metadata.add("x-user-id", firebaseUID);
  return metadata;
};
