import time
import uvicorn
import os
import torch
from fastapi import FastAPI
from fastapi import Request
from pydantic import BaseModel
from FlagEmbedding import BGEM3FlagModel
from typing import Optional, List
from datetime import datetime
import numpy as np
# pip install "optimum[onnxruntime-gpu]" transformers
from optimum.onnxruntime import ORTModelForSequenceClassification
from transformers import AutoTokenizer

app = FastAPI()

Bge_m3_model = "bge-m3"
Bge_reranker_base_onnx="bge-reranker-base-onnx-o4"
model = BGEM3FlagModel(Bge_m3_model, use_fp16=True) 
Rerank_tokenizer = AutoTokenizer.from_pretrained(Bge_reranker_base_onnx)
Rerank_model = ORTModelForSequenceClassification.from_pretrained(Bge_reranker_base_onnx)

def compute_embedding(inputs:List[str])->list:
    result = []
    embeddings = model.encode(inputs)["dense_vecs"]
    if embeddings is not None:
        data = embeddings.tolist()
        # if isinstance(data, list):
        for i, passage in enumerate(inputs):
            similarity_dict = {"index":i,"embedding": data[i]}
            result.append(similarity_dict)
    return result


def compute_similarity(query:list,inputs: List[List[float]]):
    result = []
    if len(inputs)>0:
        result = np.array(query)@np.array(inputs).T
    if result is not None:
        return result.tolist()
    return result


def compute_score_by_m3(passages,pairs: List[List[str]]):
    if len(pairs) > 0:
        result = []
        score = model.compute_score(pairs)
        scores = score["colbert"]
        if isinstance(scores, list):
            for i, passage in enumerate(passages):
                similarity_dict = {"index":i,"relevance_score": scores[i]}
                result.append(similarity_dict)
        else:
            for i, passage in enumerate(passages):
                similarity_dict = {"index":i,"relevance_score": scores[i]}
                result.append(similarity_dict)
        return result
    else:
        return None


def compute_score_by_base_onnx(passages,pairs: List[List[str]]):
    if len(pairs) > 0:
        result = []
        with torch.no_grad():
            inputs = Rerank_tokenizer(pairs, padding=True, truncation=True, return_tensors='pt', max_length=256)
            scores_data = Rerank_model(**inputs, return_dict=True).logits.view(-1, ).float()
            scores = scores_data.tolist()
            if isinstance(scores, list):
                for i, passage in enumerate(passages):
                    similarity_dict = {"index":i,"relevance_score": scores[i]}
                    result.append(similarity_dict)
            else:
                for i, passage in enumerate(passages):
                    similarity_dict = {"index":i,"relevance_score": scores[i]}
                    result.append(similarity_dict)
        return result
    else:
        return None


class RerankQuerySuite(BaseModel):
    model:str
    query: str
    passages: List[str]

@ app.post("/v1/rerank")
def api_rerank(query_suite:RerankQuerySuite,request:Request):
    result = []
    model = query_suite.model
    passages = query_suite.passages
    pairs = [[query_suite.query, passage] for passage in query_suite.passages]
    start_time = time.time()
    if model == Bge_m3_model:
        result = compute_score_by_m3(passages,pairs)
    elif model == Bge_reranker_base_onnx:
        result = compute_score_by_base_onnx(passages,pairs)
    print_log(request.url.path,start_time)
    return {"results":result}

class SimilarityQuerySuite(BaseModel):
    model:str
    query: List[float]
    input: list
@ app.post("/v1/similarity")
def api_similarity(query_suite:SimilarityQuerySuite,request:Request):
    result = []
    passages = query_suite.input
    query = query_suite.query
    start_time = time.time()
    result = compute_similarity(query,passages)
    print_log(request.url.path,start_time)
    return {"data":result}

class QuerySuite(BaseModel):
    model: str
    input: List[str]
@ app.post("/v1/embeddings")
def api_embeddings(query_suite:QuerySuite,request:Request):
    result = []
    passages = query_suite.input
    start_time = time.time()
    result = compute_embedding(passages)
    print_log(request.url.path,start_time)
    return {"data":result}

def print_log(path,start_time):
    end_time = time.time()
    print(datetime.now(),path,":spend: {:.2f}s".format(end_time - start_time))

if __name__ == '__main__':
    port = os.getenv("BGE_SERVICE_PORT")
    token = os.getenv("ACCESS_TOKEN")
    try:
        uvicorn.run(app, host='0.0.0.0', port=port,log_config="uvicorn_config.json")
    except Exception as e:
        print(f"api server start fail\n{e}")
    