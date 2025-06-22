const collections = {
  Profile: "profiles",
  Character: "characters",
};

module.exports = {
  async up(db, client) {
    const session = client.startSession();
    try {
      await session.withTransaction(async () => {
        await db.collection(collections.Profile).updateMany(
          {},
          [
            {
              $set: {
                limited_character_number: 5,
                updated_at: {
                  $ifNull: ["$created_at", new Date()],
                },
              },
            },
          ],
          { session }
        );

        await db.collection(collections.Character).updateMany(
          {},
          {
            $set: {
              created_at: new Date(),
              updated_at: new Date(),
            },
          },
          { session }
        );
      });
    } finally {
      await session.endSession();
    }
  },

  async down(db, client) {
    const session = client.startSession();
    try {
      await session.withTransaction(async () => {
        await db.collection(collections.Profile).updateMany(
          {},
          {
            $unset: {
              limited_character_number: "",
              updated_at: "",
            },
          },
          { session }
        );

        await db.collection(collections.Character).updateMany(
          {},
          {
            $unset: {
              created_at: "",
              updated_at: "",
            },
          },
          { session }
        );
      });
    } finally {
      await session.endSession();
    }
  },
};
