// Resolves which API origin to call.
//
// In production the browser and the SSR process must not use the same value:
// the browser calls the public path (/api) which Caddy proxies, while SSR runs
// inside the Docker network and can reach the API container directly — avoiding
// a pointless round trip out through the proxy and back in. In development both
// resolve to the same localhost origin.
export const useApiBaseUrl = (): string => {
  const runtimeConfig = useRuntimeConfig()

  return import.meta.server
    ? runtimeConfig.apiInternalBase
    : runtimeConfig.public.apiBase
}
