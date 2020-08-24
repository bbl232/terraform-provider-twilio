resource "random_string" "random" {
  length  = 16
  special = false
}

resource "twilio_flex_flow" "flow" {
  friendly_name    = "twilio-test-${random_string.random.result}"
  chat_service_sid = var.chat_service_sid
  channel_type     = "web"
  integration_type = "external"
  integration {
    url = "https://test.com/external"
  }
}

resource "twilio_flex_channel" "channel" {
  chat_friendly_name      = "twilio-test-${random_string.random.result}"
  chat_user_friendly_name = "twilio-test-${random_string.random.result}"
  flex_flow_sid           = twilio_flex_flow.flow.sid
  identity                = "test"
}
