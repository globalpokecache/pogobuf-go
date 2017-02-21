package client

import (
	"context"
	"github.com/globalpokecache/POGOProtos-go"
	"time"
)

var tutorialRequirements = []int32{0, 1, 3, 4, 7}

func (c *Instance) completeTutorial(ctx context.Context, tutorialState []protos.TutorialState, account string) error {
	completed := 0
	tuto := map[int32]bool{}
	for _, t := range tutorialState {
		for _, req := range tutorialRequirements {
			if req == int32(t) {
				tuto[req] = true
				completed++
			}
		}
	}

	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

	if completed == 5 {
		getPlayerProfile, err := c.GetPlayerProfileRequest("")
		if err != nil {
			return err
		}
		var requests []*protos.Request
		requests = append(requests, getPlayerProfile)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		registerBackground, err := c.RegisterBackgroundDeviceRequest("", "apple_watch")
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, registerBackground)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		return nil
	}

	if _, ok := tuto[0]; !ok {
		time.Sleep(time.Duration(2+randInt(3)) * time.Second)
		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{0}, false, false)
		if err != nil {
			return err
		}
		requests := []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[1]; !ok {
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
		requests := []*protos.Request{}
		requests = append(requests, setAvatar)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(1+randInt(1)) * time.Second)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{1}, false, false)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	getPlayerProfile, err := c.GetPlayerProfileRequest("")
	if err != nil {
		return err
	}
	var requests []*protos.Request
	requests = append(requests, getPlayerProfile)
	requests = append(requests, c.BuildCommon()...)
	requests = append(requests, getBuddyWalkedReq)
	_, err = c.Call(ctx, requests...)
	if err != nil {
		return err
	}

	registerBackground, err := c.RegisterBackgroundDeviceRequest("", "apple_watch")
	if err != nil {
		return err
	}
	requests = []*protos.Request{}
	requests = append(requests, registerBackground)
	requests = append(requests, c.BuildCommon()...)
	_, err = c.Call(ctx, requests...)
	if err != nil {
		return err
	}

	if _, ok := tuto[3]; !ok {
		getDownloadsURLs, err := c.GetDownloadURLsRequest([]string{
			"1a3c2816-65fa-4b97-90eb-0b301c064b7a/1477084786906000",
			"e89109b0-9a54-40fe-8431-12f7826c8194/1477084802881000",
		})
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, getDownloadsURLs)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(7+randInt(3)) * time.Second)
		crea := []int32{1, 4, 7}[randInt(3)]

		encounterRequest, err := c.EncounterTutorialCompleteRequest(crea)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, encounterRequest)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		getPlayerRequest, err := c.GetPlayerRequest("US", "en", "America/Chicago")
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, getPlayerRequest)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[4]; !ok {
		time.Sleep(time.Duration(5+randInt(7)) * time.Second)

		claimCodename, err := c.ClaimCodenameRequest(account)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, claimCodename)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{4}, false, false)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	if _, ok := tuto[7]; !ok {
		time.Sleep(time.Duration(4+randInt(3)) * time.Second)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{7}, false, false)
		if err != nil {
			return err
		}
		requests = []*protos.Request{}
		requests = append(requests, markComplete)
		requests = append(requests, c.BuildCommon()...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return err
		}
	}

	return nil
}
