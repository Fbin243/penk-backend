import "./bootstrap";

import express from "express";
import { Webhooks } from "@polar-sh/express";
import { PaymentModel } from "./utils/database/mongo";

const app = express();
const PORT = 8095;

// TODO: Read these values from the database
const proTierMonthlyCredit = 100;
const costPerCredit = 0.1;

app.use(express.json()).post(
    "/polar/webhooks",
    Webhooks({
        webhookSecret: process.env.POLAR_WEBHOOK_SECRET!,
        onOrderPaid: async (payload) => {
            if (payload.data.subscription) {
                // Handle subscription payments
                const email = payload.data.customer.email;
                const monthlyCredit = proTierMonthlyCredit;
                const periodEnd = payload.data.subscription.currentPeriodEnd;
                PaymentModel.findOneAndUpdate(
                    { email },
                    {
                        $set: {
                            monthly_credit: monthlyCredit,
                            period_end: periodEnd,
                        },
                    },
                    { new: true, upsert: true }
                ).catch((err) => {
                    console.error("Error updating payment model:", err);
                });
            } else {
                // Handle one-time credit purchases
                const email = payload.data.customer.email;
                const creditsPurchased =
                    payload.data.subtotalAmount / 100 / costPerCredit; // The amount is in cents
                PaymentModel.findOneAndUpdate(
                    { email },
                    {
                        $inc: {
                            persistent_credit: creditsPurchased,
                        },
                    },
                    { new: true, upsert: true }
                ).catch((err) => {
                    console.error("Error updating payment model:", err);
                });
            }
        },
    })
);

app.listen(PORT, () => {
    console.log(`Webhook server running at http://localhost:${PORT}`);
});
