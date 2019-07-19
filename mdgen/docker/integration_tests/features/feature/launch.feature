Feature: The java app

  Scenario: Launching the java app
    Given the java app exists
    And the login credentials exist
    When I run the java executable with no parameters
    Then I see that the application complains about missing parameters

  Scenario: Launching the java app with decent parameters
    Given the java app exists
    And the login credentials exist
    When I run the java executable with some appropriate parameters
    Then I see that the application complains about missing parameters
