// Temporary in-memory store for citizen photos
// One-time use, 5 minute TTL

const imageStore = new Map<string, { data: string; expires: number }>();

// Cleanup expired entries periodically
if (typeof setInterval !== "undefined") {
  setInterval(() => {
    const now = Date.now();
    imageStore.forEach((val, key) => {
      if (now > val.expires) imageStore.delete(key);
    });
  }, 60_000);
}

export { imageStore };
