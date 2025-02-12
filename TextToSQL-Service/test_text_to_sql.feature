Feature: Text-to-SQL API Testing
  Testing the endpoints of the Text-to-SQL API.

  Scenario: Check health endpoint
    Given the Text-to-SQL service is running
    When I send a GET request to "/health"
    Then the response status code should be 200
    And the response body should be valid JSON with a healthy status

  Scenario: Translate query to SQL
    When I send a POST request to "/text-to-sql/simple" with a valid query
    Then the POST response status code should be 200
