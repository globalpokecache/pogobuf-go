package client

import (
	"context"

	"github.com/globalpokecache/POGOProtos-go"
)

var tutorialRequirements = []protos.TutorialState{
	protos.TutorialState_LEGAL_SCREEN,                   // 0
	protos.TutorialState_AVATAR_SELECTION,               // 1
	protos.TutorialState_POKEMON_CAPTURE,                // 3
	protos.TutorialState_NAME_SELECTION,                 // 4
	protos.TutorialState_FIRST_TIME_EXPERIENCE_COMPLETE, // 7
}

func (c *Instance) completeTutorial(ctx context.Context, tutorialState []protos.TutorialState, account string, assets []string) (bool, error) {
	getBuddyWalkedReq, _ := c.GetBuddyWalkedRequest()

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
		return true, nil
	}

	if _, ok := tuto[protos.TutorialState_LEGAL_SCREEN]; !ok {
		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_LEGAL_SCREEN}, false, false)
		if err != nil {
			return false, err
		}
		requests := []*protos.Request{markComplete}
		requests = append(requests, c.BuildCommon(false)...)
		_, err = c.Call(ctx, markComplete)
		if err != nil {
			return false, err
		}

		randSleep(350, 525)

		getPlayerRequest, err := c.GetPlayerRequest("US", "en", "America/Chicago")
		if err != nil {
			return false, err
		}
		_, err = c.Call(ctx, getPlayerRequest)
		if err != nil {
			return false, err
		}
		randSleep(1000, 1100)
	}

	if _, ok := tuto[protos.TutorialState_AVATAR_SELECTION]; !ok {
		randSleep(5000, 5100)
		listAvatar, err := c.ListAvatarCustomizationsRequest(0, []protos.Slot{}, []protos.Filter{2})
		if err != nil {
			return false, err
		}
		requests := []*protos.Request{listAvatar}
		requests = append(requests, c.BuildCommon(false)...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

		randSleep(7000, 14000)
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
			return false, err
		}
		requests = []*protos.Request{setAvatar}
		requests = append(requests, c.BuildCommon(false)...)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

		randSleep(500, 4000)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_AVATAR_SELECTION}, false, false)
		if err != nil {
			return false, err
		}
		_, err = c.Call(ctx, markComplete)
		if err != nil {
			return false, err
		}

		randSleep(500, 1000)
	}

	if _, ok := tuto[protos.TutorialState_POKEMON_CAPTURE]; !ok {
		randSleep(700, 900)

		getDownloadsURLs, err := c.GetDownloadURLsRequest(assets)
		if err != nil {
			return false, err
		}
		requests := []*protos.Request{getDownloadsURLs}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

		randSleep(7000, 10300)

		crea := []int32{1, 4, 7}[randInt(3)]

		encounterRequest, err := c.EncounterTutorialCompleteRequest(crea)
		if err != nil {
			return false, err
		}
		requests = []*protos.Request{encounterRequest}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

		randSleep(400, 500)

		getPlayerRequest, err := c.GetPlayerRequest("US", "en", "America/Chicago")
		if err != nil {
			return false, err
		}
		requests = []*protos.Request{getPlayerRequest}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}
	}

	if _, ok := tuto[protos.TutorialState_NAME_SELECTION]; !ok {
		randSleep(12000, 18000)

		claimCodename, err := c.ClaimCodenameRequest(account)
		if err != nil {
			return false, err
		}
		requests := []*protos.Request{claimCodename}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

		randSleep(700, 800)

		getPlayerRequest, err := c.GetPlayerRequest("US", "en", "America/Chicago")
		if err != nil {
			return false, err
		}
		requests = []*protos.Request{getPlayerRequest}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

		randSleep(130, 200)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_NAME_SELECTION}, false, false)
		if err != nil {
			return false, err
		}
		requests = []*protos.Request{markComplete}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}

	}

	if _, ok := tuto[protos.TutorialState_FIRST_TIME_EXPERIENCE_COMPLETE]; !ok {
		randSleep(3900, 4500)

		markComplete, err := c.MarkTutorialCompleteRequest([]protos.TutorialState{protos.TutorialState_FIRST_TIME_EXPERIENCE_COMPLETE}, false, false)
		if err != nil {
			return false, err
		}
		requests := []*protos.Request{markComplete}
		requests = append(requests, c.BuildCommon(false)...)
		requests = append(requests, getBuddyWalkedReq)
		_, err = c.Call(ctx, requests...)
		if err != nil {
			return false, err
		}
	}

	//    if starter_id != nil {
	//     await self.random_sleep(4, 5)
	//     request = self.api.create_request()
	//     request.set_buddy_pokemon(pokemon_id=starter_id)
	//     await self.call(request, action=2)
	//     await self.random_sleep(.8, 1.2)
	//    }

	randSleep(200, 300)

	return false, nil
}
