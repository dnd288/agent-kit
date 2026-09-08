# Diagnosis

## Root cause

<!-- The specific line, query, condition or assumption that is wrong. -->

## The observation that proves it

<!-- The log line, failing assertion, query plan, network trace or diff that
     distinguishes this cause from the alternatives. A diagnosis with no evidence is
     a hypothesis, and shipping a fix for a hypothesis is how a defect comes back. -->

## Why it was not caught

<!-- Which test should have failed and did not, or which check does not exist. This
     is the part that stops the same class recurring, and it usually points at the
     test to write. -->

## What else depends on this

<!-- The blast radius of the FIX, not of the defect. Every caller of the thing being
     changed, and every consumer that is relying on the behaviour as it stands today.

     Name the backward-compatibility question and answer it: a response schema field,
     a contract enum, a migration running while the previous release is still live, a
     share link already sent, a session cookie's name or path, a shared e2e step whose
     sentence other features also say. "Nothing else calls it" is a fine answer once
     it has been looked for.

     This section chooses the run set in verification.md. A fix whose impact was never
     mapped is verified against the symptom and nothing else. -->

## Alternatives ruled out

<!-- What else could have produced this symptom, and what ruled each out. -->
