Feature: Depositing cash
Rule: Deposited amounts must be strictly positive

Scenario: Successful deposit into account
  Given Bob tiene 100.00 en su cuenta
  When Bob intenta abonar 50.00
  Then el abono se realiza con éxito y el nuevo saldo es 150.00

Scenario: Declined deposit with zero or negative amount
  Given Bob tiene 100.00 en su cuenta
  When Bob intenta abonar -20.00
  Then el abono es rechazado por monto inválido
