/**
 * The API origin, and helpers that build URLs against it.
 *
 * One module rather than a literal per spec: the port is configurable, and a second copy
 * carrying a different default is a second place to change — which is the way two specs
 * start talking to two different services and neither of them is wrong.
 */

const apiPort = process.env.E2E_API_PORT ?? '3001';

/**
 * The API's origin. Used by `seed.setup.ts` to sign in and by any step that talks to the
 * API directly — the half of a flow that has no reachable surface yet.
 *
 * Replace the port and the path prefix with your project's actual values.
 */
export const API_ORIGIN = `http://localhost:${apiPort}`;

/**
 * The API's health endpoint.
 *
 * Outside the application prefix on purpose: the one route that answers regardless of what
 * any base path is set to.
 */
export function apiHealthUrl(): string {
  return `${API_ORIGIN}/health`;
}
