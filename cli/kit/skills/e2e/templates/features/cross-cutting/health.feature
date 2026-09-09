@cross-cutting
Feature: The application answers

  The cheapest claims in the suite, and the ones every other scenario silently depends on: that
  the front door is where a visitor with no session lands, that the sign-in screen renders from
  its own route, and that both halves of the application say they are alive.

  They are worth their own file because their failures are indistinguishable from everybody
  else's at a glance. A run where every scenario fails at "signed in" is either a broken product
  or an application that never started, and these four say which within seconds.

  Rule: A visitor with no session meets the front door, not an error

    Scenario: The application root sends a visitor with no session to Sign in
      Given nobody is signed in
      When they open the application root
      Then they are asked to sign in, rather than shown an error

    Scenario: The sign-in screen renders from its own route
      Given nobody is signed in
      When they open the sign-in screen
      Then it offers an email, a password and a way in

  Rule: Both halves of the application answer for themselves

    Scenario: The API says it is healthy
      When the API's health endpoint is asked
      Then it answers that it is ok

    Scenario: The client says it is healthy
      When the client's health endpoint is asked
      Then it answers that it is ok
