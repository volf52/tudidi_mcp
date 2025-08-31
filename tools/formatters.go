package tools

import (
	"fmt"
	"strings"
	"tudidi_mcp/tudidi"
)

// FormatProjectsText formats a slice of projects into readable text
func FormatProjectsText(projects []tudidi.Project, prefix string) string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("%s:\n\n", prefix))

	for _, project := range projects {
		text.WriteString(formatSingleProject(project))
		text.WriteString("---\n\n")
	}

	return text.String()
}

// FormatTasksText formats a slice of tasks into readable text
func FormatTasksText(tasks []tudidi.Task) string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("Found %d tasks:\n\n", len(tasks)))

	for _, task := range tasks {
		text.WriteString(formatSingleTask(task))
		text.WriteString("---\n\n")
	}

	return text.String()
}

// formatSingleProject formats one project with relevant fields
func formatSingleProject(project tudidi.Project) string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("ID: %d\n", project.ID))
	text.WriteString(fmt.Sprintf("Name: %s\n", project.Name))
	if project.Description != "" {
		text.WriteString(fmt.Sprintf("Description: %s\n", project.Description))
	}
	if project.Priority != "" {
		text.WriteString(fmt.Sprintf("Priority: %s\n", project.Priority))
	}
	text.WriteString(fmt.Sprintf("Active: %t\n", project.Active))
	if project.DueDateAt != "" {
		text.WriteString(fmt.Sprintf("Due Date: %s\n", project.DueDateAt))
	}
	return text.String()
}

// formatSingleTask formats one task with relevant fields
func formatSingleTask(task tudidi.Task) string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("ID: %d\n", task.ID))
	text.WriteString(fmt.Sprintf("Name: %s\n", task.Name))
	if task.Note != "" {
		text.WriteString(fmt.Sprintf("Note: %s\n", task.Note))
	}
	text.WriteString(fmt.Sprintf("Status: %d\n", task.Status))
	text.WriteString(fmt.Sprintf("Priority: %d\n", task.Priority))
	if task.DueDate != "" {
		text.WriteString(fmt.Sprintf("Due Date: %s\n", task.DueDate))
	}
	if task.ProjectID != 0 {
		text.WriteString(fmt.Sprintf("Project ID: %d\n", task.ProjectID))
	}
	text.WriteString(fmt.Sprintf("Today: %t\n", task.Today))
	if task.CompletedAt != "" {
		text.WriteString(fmt.Sprintf("Completed: %s\n", task.CompletedAt))
	}
	return text.String()
}
