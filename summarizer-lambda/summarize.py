import os
import boto3
import google.generativeai as genai
import supabase


def store_response(content, key):
    client = supabase.create_client(
        os.getenv("SUPABASE_URL"), os.getenv("SUPABASE_KEY")
    )
    response = ""
    try:
        response = (
            client.table(os.getenv("TABLE_NAME"))
            .insert({"id": key, "response": content})
            .execute()
        )
        response = "summary stored successfully"
    except Exception as err:
        response = f"error: {err.message.split("\\")[0]}"
    finally:
        return response


def summarizer_handler(event, context):
    genai.configure(api_key=os.getenv("GEMINI_API_KEY"))
    model = genai.GenerativeModel("gemini-1.5-flash")
    s3_client = boto3.client("s3")
    boto3.set_stream_logger("botocore", level="DEBUG")
    object_name = event["Records"][0]["s3"]["object"]["key"]

    response = s3_client.get_object(Bucket=os.getenv("BUCKET_NAME"), Key=object_name)
    content = response["Body"].read()
    prompt = "summarize this file"
    response = model.generate_content(
        [{"mime_type": "application/pdf", "data": content}, prompt]
    )
    res = store_response(response.text, object_name)

    if res == "summary stored successfully":
        return {"statusCode": 200, "body": res}
    return {"statusCode": 500, "body": res}
