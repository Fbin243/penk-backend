import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";

import type { ProtoGrpcType as CoreGrpcType } from "./proto/core";

const coreService = grpc.loadPackageDefinition(
  protoLoader.loadSync("../../pkg/proto/core.proto", {}),
) as unknown as CoreGrpcType;

const client = new coreService.pb.Core(
  process.env.CORE_SERVICE_ADDRESS || "localhost:50051",
  grpc.credentials.createInsecure(),
);

client.CheckPermission(
  {
    profileID: "",
    characterID: "",
    categoryID: "",
  },
  (err, res) => {
    if (err) {
      console.error(err);
    } else {
      console.log("--> result", res);
    }
  },
);
