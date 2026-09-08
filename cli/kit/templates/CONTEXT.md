# Context glossary

A glossary keeps people and agents from using the same word for different concepts. Add terms here when a domain word collides with everyday language, another technical domain, a vendor name, or a nearby product concept.

Use this file before inventing names. If two terms feel similar, write the distinction here before encoding it in code.

## How to add an entry

Use this format:

```md
### Term

Definition: One or two sentences defining the term in this project.

Use when: The situations where this term is the correct one.

Do not use for: Similar concepts that need another word.

Related: Links to specs, ADRs, code modules, or external standards.
```

## Example entries

### Account

Definition: A login identity that can authenticate to the product.

Use when: Describing sign-in, sessions, profile settings, and ownership of actions.

Do not use for: A customer organization, billing entity, tenant, or workspace unless the project explicitly treats those as the same concept.

Related: `docs/engineering/authentication.md`, `docs/adr/0001-identity-model.md`

### Workspace

Definition: A shared area where a group collaborates on project resources.

Use when: Describing membership, shared settings, invitations, and resource grouping.

Do not use for: A local package directory in a monorepo. Use "package" or "project workspace" for that if needed.

Related: `docs/product/workspaces.md`

### Theme

Definition: A named set of visual tokens that changes the application's appearance.

Use when: Describing colour, typography, radius, shadow, spacing, and dark/light variants.

Do not use for: A user-selected content style, category, template, or mood.

Related: `docs/engineering/design-system.md`

### Token

Definition: A design-system value or machine-readable credential, depending on context. Always qualify the word when both meanings exist in the project.

Use when: Writing "design token", "access token", "refresh token", or "API token".

Do not use for: An unspecified string with special powers.

Related: `docs/engineering/security.md`, `docs/engineering/design-system.md`

## Common disambiguation patterns

- If a word names both a person and a software component, rename one in code or always qualify it.
- If a word names both a customer concept and an internal implementation detail, reserve the plain word for the customer concept.
- If a vendor uses a term differently from the product, write both definitions and cite the vendor context.
- If a term changed during the project, record the old term as deprecated and point to the replacement.
- If two teams use different names for the same thing, pick the canonical product term and list aliases.
- If a term appears in URLs, database tables, or API fields, document whether it can be renamed and what migration would be required.

## Project terms

Add real project entries below this line.
