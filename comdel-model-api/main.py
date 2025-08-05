from fastapi import FastAPI, Query
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import torch

app = FastAPI()

model = AutoModelForSequenceClassification.from_pretrained("KeyzarRasya/yt-com")
tokenizer = AutoTokenizer.from_pretrained("KeyzarRasya/yt-com")
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
model = model.to(device)
model.eval()

def predict(text):
    inputs = tokenizer(text, return_tensors="pt", truncation=True, padding=True).to(device)

    with torch.no_grad():
        outputs = model(**inputs)

    logits = outputs.logits
    predicted_class = torch.argmax(logits, dim=1).item()

    return predicted_class

@app.get("/comment/detect")
def detect_comment(comment: str = Query(default="")):
    print(comment)
    predicted = predict(comment)
    return {"result":predicted}
