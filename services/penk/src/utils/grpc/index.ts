import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";

import type { ProtoGrpcType as CoreGrpcType } from "./proto/core_service";
import type { ProtoGrpcType as TimeTrackingGrpcType } from "./proto/timetracking_service";

const coreService = grpc.loadPackageDefinition(
  protoLoader.loadSync("../../proto/core/core_service.proto", {
    includeDirs: ["../../proto"],
  }),
) as unknown as CoreGrpcType;

export const coreClient = new coreService.core.Core(
  `localhost:${process.env.CORE_GRPC_PORT || "50051"}`,
  grpc.credentials.createInsecure(),
);

const timeTrackingService = grpc.loadPackageDefinition(
  protoLoader.loadSync("../../proto/timetracking/timetracking_service.proto", {
    includeDirs: ["../../proto"],
  }),
) as unknown as TimeTrackingGrpcType;

export const timeTrackingClient = new timeTrackingService.timetracking.TimeTrackingService(
  `localhost:${process.env.TIMETRACKING_GRPC_PORT || "50053"}`,
  grpc.credentials.createInsecure(),
);
