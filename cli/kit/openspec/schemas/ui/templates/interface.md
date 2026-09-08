# Interface

## Props this change exposes

<!-- The prop types of each component or screen. These ARE the specification of what
     the API must eventually return — write them from the design, not from a guessed
     schema. -->

## Contract shapes assumed

<!-- Types from the API client package this change's caller will map into the props above.
     Name them even though nothing here imports them: the UI package may not import
     the shared contract package, and stating the assumption is how the be change knows what to
     build. Write "none" if the component is presentational data-free. -->

## Store and state

<!-- What this change reads from the session store package, if anything, and what it deliberately
     does not. A composite reads nothing: data arrives as props (the state-management decision record). If this
     change needs a store, say which slice and why props would not do. -->
