import json
import requests
from pytest_bdd import scenarios, given, when, then

scenarios("test_text_to_sql.feature")

BASE_URL = "http://127.0.0.1:5006"

@given("the Text-to-SQL service is running")
def service_running():
    try:
        response = requests.get(f"{BASE_URL}/health")
        assert response.status_code == 200
    except requests.exceptions.RequestException as e:
        raise AssertionError(f"Сервис недоступен: {e}")

@when('I send a GET request to "/health"')
def send_get_request_to_health():
    global health_response
    health_response = requests.get(f"{BASE_URL}/health")

@then("the response status code should be 200")
def check_health_status_code():
    assert health_response.status_code == 200

@then("the response body should be valid JSON with a healthy status")
def check_health_response_body():
    expected_body = {"status": "healthy"}
    actual_body = health_response.json()
    print(actual_body)
    assert actual_body == expected_body

@when('I send a POST request to "/text-to-sql/simple" with a valid query')
def send_post_request_to_simple():
    global sql_response
    request_body = {"query": "Find all users that are elder than 25"}
    sql_response = requests.post(f"{BASE_URL}/text-to-sql/simple", json=request_body)

@then("the POST response status code should be 200")
def check_post_status_code():
    assert sql_response.status_code == 200
