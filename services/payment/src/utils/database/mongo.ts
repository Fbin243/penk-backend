import mongoose from "mongoose";

const conn = mongoose.createConnection(
    `mongodb+srv://${process.env.MONGO_USER}:${process.env.MONGO_PASSWORD}@${process.env.MONGO_ADDRESS}/${process.env.MONGO_DATABASE_NAME}?retryWrites=true&w=majority`,
    {
        autoIndex: true,
    }
);

const Schema = mongoose.Schema;

const ProfileSchema = new Schema({
    name: { type: String, required: true },
    email: { type: String, required: true },
    current_character_id: { type: Schema.Types.ObjectId, required: true },
});

export const ProfileModel = conn.model("profiles", ProfileSchema);

const PaymentSchema = new Schema({
    // I don't use `user_id` because a user may make a purchase before creating an account
    email: { type: String, required: true },
    monthly_credit: { type: Number, default: 0 },
    persistent_credit: { type: Number, default: 0 },
    period_end: { type: Date, default: null },
});

export const PaymentModel = conn.model("payments", PaymentSchema);
