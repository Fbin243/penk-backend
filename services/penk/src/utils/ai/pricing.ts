import { CompletionUsage } from "openai/resources/completions.mjs";
import { encoding_for_model } from "tiktoken";

/* https://platform.openai.com/docs/pricing */
export interface CompletionPricingModel {
  pricePer1MPromptTextTokens: number;
  pricePer1MPromptTextTokensCached: number;
  pricePer1MCompletionTextTokens: number;
  pricePer1MPromptAudioTokens: number;
  pricePer1MCompletionAudioTokens: number;
}

export const gpt4oMiniPricingModel: CompletionPricingModel = {
  pricePer1MPromptTextTokens: 0.15,
  pricePer1MPromptTextTokensCached: 0.075,
  pricePer1MCompletionTextTokens: 0.6,
  pricePer1MPromptAudioTokens: 10,
  pricePer1MCompletionAudioTokens: 20,
};

export const calculateCompletionUsage = (
  usage: CompletionUsage,
  pricingModel: CompletionPricingModel,
) => {
  const textInputTokens =
    usage.prompt_tokens - (usage.prompt_tokens_details?.cached_tokens ?? 0) || 0;
  const cachedTextTokens = usage.prompt_tokens_details?.cached_tokens || 0;
  const audioInputTokens = usage.prompt_tokens_details?.audio_tokens || 0;

  const textOutputTokens = usage.completion_tokens || 0;
  const audioOutputTokens = usage.completion_tokens_details?.audio_tokens || 0;

  const textInputCost =
    (textInputTokens * pricingModel.pricePer1MPromptTextTokens +
      cachedTextTokens * pricingModel.pricePer1MPromptTextTokensCached) /
    1000000;
  const audioInputCost = audioInputTokens / 1000000;
  const textOutputCost = (textOutputTokens * pricingModel.pricePer1MCompletionTextTokens) / 1000000;
  const audioOutputCost =
    (audioOutputTokens * pricingModel.pricePer1MCompletionAudioTokens) / 1000000;

  const totalCost = textInputCost + audioInputCost + textOutputCost + audioOutputCost;

  return totalCost;
};

export const calculateTranscriptionCost = (content: string) => {
  // We use gpt-4o-transcribe which costs $6 per 1M tokens
  const tokens = encoding_for_model("gpt-4o").encode(content).length;
  const cost = (tokens * 6) / 1000000;
  return cost;
};

export const calculateTtsCost = (duration: number) => {
  // We use gpt-4o-mini-tts which costs approximately $0.015 per minute
  const cost = (duration * 0.015) / 60;
  return cost;
};
