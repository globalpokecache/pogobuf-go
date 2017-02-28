package client

import (
	"context"
	"github.com/globalpokecache/POGOProtos-go"
	"time"
)

var tutorialRequirements = []protos.TutorialState{
	protos.TutorialState_LEGAL_SCREEN,                   // 0
	protos.TutorialState_AVATAR_SELECTION,               // 1
	protos.TutorialState_POKEMON_CAPTURE,                // 3
	protos.TutorialState_NAME_SELECTION,                 // 4
	protos.TutorialState_FIRST_TIME_EXPERIENCE_COMPLETE, // 7
}

func (c *Instance) completeTutorial(ctx context.Context, tutorialState []protos.TutorialState, account string) error {
	completed := 0
	tuto := map[protos.TutorialState]bool{}
	for _, t := range tutorialState {
		for _, req := range tutorialRequirements {
			if req == t {
				tuto[req] = true
				completed++
			}
		}
	}

	if completed == 5 {

		getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()
		getPlayerProfile, err := c.GetPlayerProfileRequest(c.options.AuthProvider.GetUsername())
		if err != nil {
			return err
		}
		requests := []*protos.Request{getPlayerProfile}
		requests = append(requests, c.BuildCommon(true)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		return nil
	}

	if _, ok := tuto[protos.TutorialState_LEGAL_SCREEN]; !ok {
		time.Sleep(time.Duration(2+randInt(3)) * time.Second)
		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_LEGAL_SCREEN}, false, false)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, markComplete)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[protos.TutorialState_AVATAR_SELECTION]; !ok {
		time.Sleep(time.Duration(8+randInt(7)) * time.Second)
		setAvatar, err := c.SetAvatarRequest(
			randInt(3),
			randInt(5),
			randInt(3),
			randInt(2),
			randInt(4),
			randInt(6),
			0,
			randInt(4),
			randInt(5),
		)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, setAvatar)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(1+randInt(1)) * time.Second)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_AVATAR_SELECTION}, false, false)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, markComplete)
		if err != nil {
			return err
		}
	}

	getPlayerProfile, err := c.GetPlayerProfileRequest(account)
	if err != nil {
		return err
	}
	_, err = c.Call(ctx, getPlayerProfile)
	if err != nil {
		return err
	}

	if _, ok := tuto[protos.TutorialState_POKEMON_CAPTURE]; !ok {
		getDownloadsURLs, err := c.GetDownloadURLsRequest([]string{
			"1a3c2816-65fa-4b97-90eb-0b301c064b7a/1477084786906000",
			"e89109b0-9a54-40fe-8431-12f7826c8194/1477084802881000",
		})
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, getDownloadsURLs)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(7+randInt(3)) * time.Second)
		crea := []int32{1, 4, 7}[randInt(3)]

		encounterRequest, err := c.EncounterTutorialCompleteRequest(crea)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, encounterRequest)
		if err != nil {
			return err
		}

		getPlayerRequest, err := c.GetPlayerRequest("US", "en", "America/Chicago")
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, getPlayerRequest)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[protos.TutorialState_NAME_SELECTION]; !ok {
		time.Sleep(time.Duration(5+randInt(7)) * time.Second)

		claimCodename, err := c.ClaimCodenameRequest(account)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, claimCodename)
		if err != nil {
			return err
		}

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_NAME_SELECTION}, false, false)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, markComplete)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[protos.TutorialState_FIRST_TIME_EXPERIENCE_COMPLETE]; !ok {
		time.Sleep(time.Duration(4+randInt(3)) * time.Second)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_FIRST_TIME_EXPERIENCE_COMPLETE}, false, false)
		if err != nil {
			return err
		}
		_, err = c.Call(ctx, markComplete)
		if err != nil {
			return err
		}
	}

	return nil
}
