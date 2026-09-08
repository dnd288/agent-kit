# Migration

## Does this change need one?

<!-- "No" is the common and correct answer. Say it explicitly: an unanswered
     migration question is how a schema change reaches review unnoticed. -->

## If yes

- **What it adds or alters:**
- **Why it is backward compatible with the previous release:**
- **What runs against a database at the previous migration:**

<!-- One migration per pull request. Never edit an applied one — a second migration
     corrects the first. Additive columns need a default or a nullable type, because
     the previous release is still writing rows while this one deploys. -->
