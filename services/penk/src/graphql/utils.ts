import { MessageModel, UserContextModel } from "../db";

export const getUserContext = async (profileId: string) => {
  let userContext = await UserContextModel.findOne({
    profileId: profileId,
  });

  if (!userContext) {
    userContext = await UserContextModel.create({
      profileId: profileId,
      locale: "vi",
      timezone: "Asia/Ho_Chi_Minh",
      context: "",
      preferences: {
        tone: "funny",
      },
    });
  }

  return {
    id: userContext._id.toString(),
    profileID: userContext.profileId.toString(),
    timezone: userContext.timezone,
    locale: userContext.locale,
    context: userContext.context || "",
    preferences: {
      tone: userContext.preferences?.tone || "",
    },
  };
};

export const getMessages = async (profileId: string) => {
  const messages = await MessageModel.find({
    profileId,
  });

  return messages;
};
