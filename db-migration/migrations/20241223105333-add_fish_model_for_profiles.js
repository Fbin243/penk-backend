module.exports = {
  async up(db, client) {
    const session = client.startSession();
    try {
      await session.withTransaction(async () => {
        const profiles = await db.collection("profiles").find().toArray();
        console.log(profiles);
        const fish = profiles.map((profile) => ({
          profile_id: profile._id,
          gold: 0,
          normal: 0,
          created_at: new Date(),
          updated_at: new Date(),
        }));

        await db.collection("fish").insertMany(fish, { session });
      });
    } finally {
      await session.endSession();
    }
  },

  async down(db, client) {
    const session = client.startSession();
    try {
      await session.withTransaction(async () => {
        await db.collection("fish").drop();
      });
    } finally {
      await session.endSession();
    }
  }
};
