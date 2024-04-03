package openai

import (
	"os"
	"testing"
	"time"
)

func TestRuns(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	client.SetBetaHeader(`assistants=v1`)

	// (A) testing runs
	if thread, err := client.CreateThread(CreateThreadOptions{}.
		SetMessages([]ThreadMessage{
			NewThreadMessage("How do I get six different numbers between 1 and 45 randomly?"),
		})); err != nil {
		t.Errorf("failed to create thread for testing runs (A): %s", err)
	} else {
		threadID := thread.ID

		if assistant, err := client.CreateAssistant(assistantsModel, CreateAssistantOptions{}.
			SetName("Math Tutor").
			SetInstructions("You are a personal math tutor. When asked a question, write and run Ruby code to answer the question.").
			SetTools([]Tool{
				NewCodeInterpreterTool(),
			})); err != nil {
			t.Errorf("failed to create assistant for testing runs: %s", err)
		} else {
			assistantID := assistant.ID

			// === CreateRun ===
			if created, err := client.CreateRun(threadID, assistantID, CreateRunOptions{}); err != nil {
				t.Errorf("failed to create run: %s", err)
			} else {
				runID := created.ID

				// === ListRuns ===
				if _, err := client.ListRuns(threadID, nil); err != nil {
					t.Errorf("failed to list runs: %s", err)
				}

				// === RetrieveRun ===
				if retrieved, err := client.RetrieveRun(threadID, runID); err != nil {
					t.Errorf("failed to retrieve run: %s", err)
				} else {
					if retrieved.ID != runID {
						t.Errorf("retrieved run's id is not equal to the created one's")
					}

					// === CancelRun ===
					if _, err := client.CancelRun(threadID, retrieved.ID); err != nil {
						t.Errorf("failed to cancel run: %s", err)
					}
				}

				// === ListRunSteps  ===
				if list, err := client.ListRunSteps(threadID, runID, nil); err != nil {
					t.Errorf("failed to list run steps: %s", err)
				} else {
					runStep := list.Data[0]

					// === RetrieveRunStep  ===
					if retrieved, err := client.RetrieveRunStep(threadID, runID, runStep.ID); err != nil {
						t.Errorf("failed to retrieve run step: %s", err)
					} else {
						if retrieved.ID != runStep.ID {
							t.Errorf("retrieved run step's id is not equal to the requested one's")
						}
					}
				}
			}

			if _, err := client.DeleteAssistant(assistant.ID); err != nil {
				t.Errorf("failed to delete assistant which was created for testing runs: %s", err)
			}
		}
	}
}

func TestRunsWithToolOutputs(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	client.SetBetaHeader(`assistants=v1`)

	// (B) testing runs that need submission of tool outputs
	if assistant, err := client.CreateAssistant(assistantsModel, CreateAssistantOptions{}.
		SetName("Weather Notifier").
		SetInstructions("You are a personal weather notifier. When asked about weathers, answer to the question based on current time and the place.").
		SetTools([]Tool{
			NewFunctionTool(ToolFunction{
				Name:        "get_weather",
				Description: "Determine weather in my location",
				Parameters: NewToolFunctionParameters().
					AddPropertyWithDescription("location", "string", "The city and state e.g. San Francisco, CA").
					AddPropertyWithEnums("unit", "string", "string", []string{"c", "f"}).
					SetRequiredParameters([]string{"location", "unit"}),
			}),
		})); err != nil {
		t.Errorf("failed to create assistant for testing runs: %s", err)
	} else {
		assistantID := assistant.ID

		// === CreateThreadAndRun ===
		if created, err := client.CreateThreadAndRun(assistantID, CreateThreadAndRunOptions{}); err != nil {
			t.Errorf("failed to create thread and run: %s", err)
		} else {
			runID := created.ID
			threadID := created.ThreadID

			time.Sleep(5 * time.Second) // wait for it to progress...

			if retrieved, err := client.RetrieveRun(threadID, runID); err != nil {
				t.Errorf("failed to retrieve run: %s", err)
			} else {
				if retrieved.Status == RunStatusRequiresAction && retrieved.RequiredAction.Type == "submit_tool_outputs" {
					output := "36.5c" // NOTE: get your local function's result with the generated arguments
					toolOutputs := []ToolOutput{}
					for _, toolCall := range retrieved.RequiredAction.SubmitToolOutputs.ToolCalls {
						toolCallID := toolCall.ID

						toolOutputs = append(toolOutputs, ToolOutput{
							ToolCallID: &toolCallID,
							Output:     &output,
						})
					}

					// === SubmitToolOutputsToRun ===
					if submitted, err := client.SubmitToolOutputs(threadID, runID, toolOutputs); err != nil {
						t.Errorf("failed to submit tool outputs: %s", err)
					} else {
						time.Sleep(5 * time.Second) // wait for it to progress...

						// === ModifyRun ===
						if modified, err := client.ModifyRun(threadID, submitted.ID, ModifyRunOptions{}.
							SetMetadata(map[string]string{
								"country": "Korea",
								"city":    "Seoul",
							})); err != nil {
							t.Errorf("failed to modify run: %s", err)
						} else {
							if len(modified.Metadata) != 2 {
								t.Errorf("modified metadata count is not 2")
							}
						}
					}
				} else {
					t.Errorf("run status is not %s: %s", RunStatusRequiresAction, retrieved.Status)
				}
			}

			// === ListRunSteps  ===
			if list, err := client.ListRunSteps(threadID, runID, nil); err != nil {
				t.Errorf("failed to list run steps: %s", err)
			} else {
				runStep := list.Data[0]

				// === RetrieveRunStep  ===
				if retrieved, err := client.RetrieveRunStep(threadID, runID, runStep.ID); err != nil {
					t.Errorf("failed to retrieve run step: %s", err)
				} else {
					if retrieved.ID != runStep.ID {
						t.Errorf("retrieved run step's id is not equal to the requested one's")
					}
				}
			}
		}

		if _, err := client.DeleteAssistant(assistant.ID); err != nil {
			t.Errorf("failed to delete assistant which was created for testing runs: %s", err)
		}
	}
}
