# Interface

## The mapping

<!-- Which the API client package DTO becomes which the UI package prop type, and at which route
     boundary. This is the seam this change owns: the agreement stops here and props
     start here. -->

## Store and state placement

<!-- Where each piece of state lives, and the the state-management decision record rule that put it there.
     Server → server-rendered route + server-side action; shareable → the URL; changes without
     the user acting → a client data-fetching cache; read widely, written rarely → Context;
     screen-local → a feature slice; otherwise props. The signed-in identity is the
     one store exception (the session-state decision record). -->

## What this change assumes exists

<!-- The ui change whose props it maps into, and the be change whose endpoints it
     calls. Name them. An fe change with neither is either premature or mis-typed. -->
