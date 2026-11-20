Feature: Transaction integrity and JSON audit logging
  Rule: Every transaction (withdrawal or deposit) must be atomic and recorded in JSON format

  Scenario: Successful atomic withdrawal logged as EXITO
    Given Alice has 234.56 in their account
    When Alice executes an atomic withdrawal of 200.00
    Then the new balance is 34.56
    And a JSON transaction log is recorded with status "EXITO"

  Scenario: Declined atomic withdrawal logged as FALLO without altering balance
    Given Hamza has 198.76 in their account
    When Hamza attempts an atomic withdrawal of 200.00
    Then the balance remains unchanged at 198.76
    And a JSON transaction log is recorded with status "FALLO"

  Scenario: Invalid amount logged as ERROR without altering balance
    Given Pedro has 300.00 in their account
    When Pedro attempts an atomic withdrawal of -20.00
    Then the balance remains unchanged at 300.00
    And a JSON transaction log is recorded with status "ERROR"

  Scenario: Successful atomic deposit logged as EXITO
    Given Bob has 100.00 in their account
    When Bob executes an atomic deposit of 50.00
    Then the new balance is 150.00
    And a JSON transaction log is recorded with status "EXITO"
