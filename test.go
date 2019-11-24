package main

import (
    "fmt"
)

type Job struct {
    Title string
    Description string
}

func (j *Job) SetTitle(title string) {
    j.Title = title
}

func (j *Job) SetDescription(description string) {
    j.Description = description
}

func (j *Job) GetTitle() string {
    return j.Title
}

func (j *Job) GetDescription() string {
    return j.Description
}

func main() {
    j := Job{}
    j.SetTitle("Programmer")
    j.SetDescription("Someone who programs..")
    fmt.Printf("My job is %s and I'm %s\n", j.GetTitle(), j.GetDescription())
}
