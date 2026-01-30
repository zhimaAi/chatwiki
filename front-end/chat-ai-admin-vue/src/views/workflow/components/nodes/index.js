
import { register } from '@logicflow/vue-node-registry'
import sessionTriggerNode from './session-trigger-node/index.js'
import startNode from './start-node/index.js'
import questionNode from './question-node/index.js'
import actionNode from './action-node/index.js'
import qaNode from './qa-node/index.js'
import aiDialogueNode from './ai-dialogue-node'
import httpNode from './http-node/index.js'
import httpToolNode from './http-tool-node/index.js'
import knowledgeBaseNode from './knowledge-base-node/index.js'
import endNode from './end-node/index.js'
import explainNode from './explain-node/index.js'
import variableAssignmentNode from './variable-assignment-node/index.js'
import judgeNode from './judge-node/index.js'
import specifyReplyNode from './specify-reply-node/index.js'
import immediatelyReplyNode from './immediately-reply-node/index.js'
import parameterExtractionNode from './parameter-extraction-node/index'
import problemOptimizationNode from './problem-optimization-node/index.js'
import addDataNode from './add-data-node/index.js'
import updateDataNode from './update-data-node/index.js'
import deleteDataNode from './delete-data-node/index.js'
import selectDataNode from './select-data-node/index.js'
import codeRunNode from './code-run-node/index.js'
import mcpNode from "./mcp-node/index.js";
import zmPluginsNode from "./zm-plugins-node/index.js";
import customGroupNode from './custom-group-node/index.js'
import groupStartNode from './group-start-node/index.js'
import terminateMode from './terminate-node/index.js'
import timingTriggerNode from './timing-trigger-node/index.js'
import officialTriggerNode from './official-trigger-node/index.js'
import batchGroupNode from './batch-group-node/index.js'
import imageGenerationNode from './image-generation-node/index.js'
import jsonNode from './json-node/index.js'
import jsonReverseNode from './json-reverse-node/index.js'
import voiceSynthesisNode from "./voice-synthesis-node/index.js";
import voiceCloneNode from "./voice-clone-node/index.js";
import webhookTriggerNode from './webhook-trigger-node/index.js'
import importLibraryNode from './import-library-node/index.js'
import zmWorkflowNode from "./zm-workflow-node/index.js";



export default function registerCustomNodes(lf) {
  register(sessionTriggerNode, lf)
  register(startNode, lf)
  register(questionNode, lf)
  register(actionNode, lf)
  register(httpNode, lf)
  register(httpToolNode, lf)
  register(qaNode, lf)
  register(aiDialogueNode, lf)
  register(knowledgeBaseNode, lf)
  register(endNode, lf)
  register(terminateMode, lf)
  register(explainNode, lf)
  register(judgeNode, lf)
  register(variableAssignmentNode, lf)
  register(specifyReplyNode, lf)
  register(immediatelyReplyNode, lf)
  register(parameterExtractionNode, lf)
  register(problemOptimizationNode, lf)
  register(addDataNode, lf)
  register(updateDataNode, lf)
  register(deleteDataNode, lf)
  register(selectDataNode, lf)
  register(codeRunNode, lf)
  register(mcpNode, lf)
  register(zmPluginsNode, lf)
  register(customGroupNode, lf)
  register(groupStartNode, lf)
  register(timingTriggerNode, lf)
  register(officialTriggerNode, lf)
  register(batchGroupNode, lf)
  register(imageGenerationNode, lf)
  register(jsonNode, lf)
  register(jsonReverseNode, lf)
  register(voiceSynthesisNode, lf)
  register(voiceCloneNode, lf)
  register(webhookTriggerNode, lf)
  register(importLibraryNode, lf)
  register(zmWorkflowNode, lf)
}
