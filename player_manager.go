package grouplay

import (
	"fmt"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

// Id as key
var players map[string]*GamePlayer

func init() {
	players = make(map[string]*GamePlayer)
}

type GamePlayer struct {
	ID              string
	Name            string
	Index           int
	Session         *sockjs.Session
	GroupHosted     *GameGroup
	GroupJoined     *GameGroup
	GroupSpectating *GameGroup
	InGame          bool
}

func (p *GamePlayer) Update(session sockjs.Session, id string, name string) {
	oldId := p.ID
	fmt.Println("old id is", oldId)
	p.ID = id
	p.Name = name
	p.Session = &session
	delete(players, oldId)
	fmt.Println("old id removed", oldId)
	players[p.ID] = p
	fmt.Println("new id added", p.ID)
}

func Register(session sockjs.Session, oldId string, name string) error {
	if player, ok := FindPlayer(oldId); ok {
		player.Update(session, session.ID(), name)
		fmt.Println("Find an existed player & update it")
		return nil
	} else {
		if p := FindPlayerByName(name); p != nil {
			return NewError("The name is already registered!")
		}
		id := session.ID()
		players[id] = &GamePlayer{id, name, 0, &session, nil, nil, nil, false}
		fmt.Println("Register as new")
		return nil
	}
}

func FindPlayerByName(name string) *GamePlayer {
	for _, p := range players {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func FindPlayer(id string) (player *GamePlayer, ok bool) {
	fmt.Println("Try to find a player with id", id)
	if id == "" {
		return player, false
	}
	player, ok = players[id]
	return player, ok
}

func (p *GamePlayer) CreateGroup(game *Game, max int, allowSpectate bool) (bool, error) {
	if p.GroupHosted != nil {
		fmt.Println("You already hosted a group")
		return false, NewError("You already hosted a group.")
	}
	fmt.Println("group hosted", p.GroupHosted)
	if p.GroupJoined != nil {
		fmt.Println("You already joined a group")
		return false, NewError("You already joined a group.")
	}
	fmt.Println("group joined", p.GroupJoined)
	group := CreateGroup(game, p, max, allowSpectate)
	p.GroupHosted = group
	fmt.Println("A group created by player", p.ID)
	if err := group.Join(p); err == nil {
		p.GroupJoined = group
		return true, nil
	} else {
		return true, err
	}
}

func (p *GamePlayer) JoinGroup(id string) (bool, error) {
	if p.GroupJoined != nil {
		fmt.Println("Already joined a group.")
		return false, NewError("You Already joined a group.")
	}
	if group, ok := FindGroup(id); ok {
		if err := group.Join(p); err == nil {
			p.GroupJoined = group
			return true, nil
		} else {
			return false, err
		}
	}
	fmt.Println("Target group not found.")
	return false, NewError("Target group not found.")
}

func (p *GamePlayer) SpectateGame(id string) (bool, error) {
	if p.GroupSpectating != nil {
		fmt.Println("Already spectate a group.")
		return false, NewError("You Already spectate a group.")
	}
	if group, ok := FindGroup(id); ok {
		if err := group.Spectate(p); err == nil {
			p.GroupSpectating = group
			return true, nil
		} else {
			return false, err
		}
	}
	fmt.Println("Target group not found.")
	return false, NewError("Target group not found.")
}

func (p *GamePlayer) ExitGroup(id string) (bool, error) {
	if p.GroupJoined != nil {
		if err := p.GroupJoined.Exit(p); err == nil {
			return true, err
		} else {
			return false, err
		}
	}
	return false, NewError("You haven't joined any group")
}
