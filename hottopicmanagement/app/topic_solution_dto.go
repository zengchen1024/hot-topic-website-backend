package app

import "fmt"

type CmdToAddTopicSolution []OptionalTopic

func (cmd CmdToAddTopicSolution) init() {
	for i := range cmd {
		cmd[i].init()
	}
}

func (cmd CmdToAddTopicSolution) Validate() error {
	for i := range cmd {
		topic := &cmd[i]
		items := topic.DiscussionSources

		for j := range items {
			resolved, unresolved := items[j].filterout()
			if len(unresolved) != 0 && len(resolved) != 1 {
				return fmt.Errorf(
					"resolved num is not 1, topic:%s, resolved id:%d",
					topic.Title, resolved[0].Id,
				)
			}
		}
	}

	return nil
}
