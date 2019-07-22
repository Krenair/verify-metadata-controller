Feature: The java app

  Background:
    Given the java app exists
    And the login credentials exist


  Scenario: Launching the java app
    When I run the java executable with no parameters
    Then I see that the application complains about missing parameters

  Scenario: Launching the java app with decent parameters
    When I run the java executable with some appropriate parameters
    Then I see that the application complains about missing parameters

  Scenario: Getting a signed dojamaflip
    When I run the java executable with some appropriate parameters
    And I request some signed something
    Then I get a signed something