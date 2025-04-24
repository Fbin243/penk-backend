export const isoStringToUnixSeconds = (isoString?: string | null): number | undefined => {
  if (!isoString) return undefined;
  return Math.floor(new Date(isoString).getTime() / 1000);
};
