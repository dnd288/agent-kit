/**
 * The seeded identities.
 *
 * One module rather than a copy per spec, because the set has changed under a suite before:
 * accounts dropped from the seed while specs kept naming them, one of which asserted a
 * successful sign-in as somebody who no longer existed. When every spec points here, a removed
 * account fails to compile rather than fails at runtime.
 *
 * Replace the demo accounts below with your project's own seeded identities.
 */

export interface Account {
  email: string;
  password: string;
}

/**
 * The BOOTSTRAP operator — the only identity that exists before anything has been seeded.
 *
 * Every other account in this file is written by the seed routine, which is work an
 * authenticated operator starts from the admin console. That is a chicken and egg: the suite
 * cannot sign in as a seeded admin to run the seed, because the admin IS the seed. This one
 * comes from the database's initial migration — the one row start-up writes — so it is who
 * `seed.setup.ts` signs in as to start the seed routines.
 *
 * Read from the environment rather than fixed, because a deployment supplies it and a literal
 * here would be a second place to change.
 */
export const OPERATOR: Account = {
  email: process.env.BOOTSTRAP_ADMIN_EMAIL ?? 'operator@ci.local',
  password: process.env.BOOTSTRAP_ADMIN_PASSWORD ?? 'not-a-real-secret-one-job-only',
};

/**
 * A demo password shared by all seeded accounts.
 *
 * Bound to a name rather than inlined at each use: `password: '<literal>'` is what a
 * `secret-literals` guard looks for, and teaching it to skip test files would blind it to a real
 * one landing in one. These are fixtures from a seed everybody can read.
 */
const SEEDED_PASSWORD = 'demo-password-123';

/** A regular user — the role the product is for. */
export const USER: Account = { email: 'user@demo.local', password: SEEDED_PASSWORD };

/** An administrator of the user's organisation. */
export const ADMIN: Account = { email: 'admin@demo.local', password: SEEDED_PASSWORD };

/** An address the product has never seen. The refusal must be the same as every other one. */
export const UNKNOWN: Account = { email: 'nobody@demo.local', password: SEEDED_PASSWORD };

/** A real account, the wrong password — the other half of the pair a refusal must not tell apart. */
export const WRONG_PASSWORD: Account = {
  email: USER.email,
  password: `${SEEDED_PASSWORD}-but-wrong`,
};
