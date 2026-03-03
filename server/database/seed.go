package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GenerateFakeData generates dummy data for bands, releases, and projects tables
func (q *Queries) GenerateFakeData(ctx context.Context) error {
	items, err := q.GetBands(ctx)
	if err != nil {
		return fmt.Errorf("failed to get bands: %w", err)
	}
	if len(items) > 0 {
		return fmt.Errorf("data already exists, aborting")
	}

	// Bands with variety of name lengths
	bands := []struct {
		name    string
		country string
	}{
		{"Metallica", "USA"},
		{"Iron Maiden", "UK"},
		{"Slayer", "USA"},
		{"Megadeth", "USA"},
		{"Black Sabbath", "UK"},
		{"Judas Priest", "UK"},
		{"Pantera", "USA"},
		{"Sepultura", "Brazil"},
		{"Death", "USA"},
		{"Morbid Angel", "USA"},
		{"Cannibal Corpse", "USA"},
		{"Obituary", "USA"},
		{"Napalm Death", "UK"},
		{"Carcass", "UK"},
		{"At the Gates", "Sweden"},
		{"Dark Tranquillity", "Sweden"},
		{"In Flames", "Sweden"},
		{"Opeth", "Sweden"},
		{"Amon Amarth", "Sweden"},
		{"Arch Enemy", "Sweden"},
		{"Children of Bodom", "Finland"},
		{"Dimmu Borgir", "Norway"},
		{"Emperor", "Norway"},
		{"Mayhem", "Norway"},
		{"Burzum", "Norway"},
		{"Darkthrone", "Norway"},
		{"Immortal", "Norway"},
		{"Satyricon", "Norway"},
		{"Gorgoroth", "Norway"},
		{"Behemoth", "Poland"},
		{"Nile", "USA"},
		{"Suffocation", "USA"},
		{"Cryptopsy", "Canada"},
		{"Gorguts", "Canada"},
		{"Meshuggah", "Sweden"},
		{"Gojira", "France"},
		{"Lamb of God", "USA"},
		{"Mastodon", "USA"},
		{"Neurosis", "USA"},
		{"Isis", "USA"},
		{"Cult of Luna", "Sweden"},
		{"The Ocean", "Germany"},
		{"Baroness", "USA"},
		{"High on Fire", "USA"},
		{"Sleep", "USA"},
		{"Electric Wizard", "UK"},
		{"Cathedral", "UK"},
		{"My Dying Bride", "UK"},
		{"Paradise Lost", "UK"},
		{"Anathema", "UK"},
		{"Cradle of Filth", "UK"},
		{"Dimmu Borgir and the Symphony of Eternal Darkness", "Norway"},
		{"The Berzerker Presents: Extreme Audio Warfare", "Australia"},
		{"Wolves in the Throne Room", "USA"},
		{"Agalloch", "USA"},
	}

	bandIDs := make(map[string]string) // map band name to ID
	for _, band := range bands {
		id := uuid.NewString()
		bandIDs[band.name] = id
		_, err := q.CreateBand(ctx, CreateBandParams{
			ID:      id,
			Name:    band.name,
			Country: band.country,
		})
		if err != nil {
			return fmt.Errorf("failed to create band %s: %w", band.name, err)
		}
	}

	// Releases with variety of types and name lengths
	releases := []struct {
		bandName string
		name     string
		rType    string
		number   string
	}{
		// Short names
		{"Metallica", "Kill 'Em All", "album", "1"},
		{"Metallica", "Ride the Lightning", "album", "2"},
		{"Metallica", "Master of Puppets", "album", "3"},
		{"Iron Maiden", "The Number of the Beast", "album", "4"},
		{"Iron Maiden", "Powerslave", "album", "5"},
		{"Slayer", "Reign in Blood", "album", "6"},
		{"Slayer", "South of Heaven", "album", "7"},
		{"Megadeth", "Peace Sells... but Who's Buying?", "album", "8"},
		{"Megadeth", "Rust in Peace", "album", "9"},
		{"Black Sabbath", "Paranoid", "album", "10"},
		{"Black Sabbath", "Master of Reality", "album", "11"},
		{"Judas Priest", "British Steel", "album", "12"},
		{"Judas Priest", "Painkiller", "album", "13"},
		{"Pantera", "Cowboys from Hell", "album", "14"},
		{"Pantera", "Vulgar Display of Power", "album", "15"},
		{"Sepultura", "Beneath the Remains", "album", "16"},
		{"Sepultura", "Arise", "album", "17"},
		{"Death", "Scream Bloody Gore", "album", "18"},
		{"Death", "Leprosy", "album", "19"},
		{"Morbid Angel", "Altars of Madness", "album", "20"},
		{"Cannibal Corpse", "Tomb of the Mutilated", "album", "21"},
		{"Obituary", "The End Complete", "album", "22"},
		{"Napalm Death", "Scum", "album", "23"},
		{"Carcass", "Symphonies of Sickness", "album", "24"},
		{"At the Gates", "Slaughter of the Soul", "album", "25"},
		{"Dark Tranquillity", "The Gallery", "album", "26"},
		{"In Flames", "The Jester Race", "album", "27"},
		{"Opeth", "Blackwater Park", "album", "28"},
		{"Opeth", "Ghost Reveries", "album", "29"},
		{"Amon Amarth", "Twilight of the Thunder God", "album", "30"},
		{"Arch Enemy", "Wages of Sin", "album", "31"},
		{"Children of Bodom", "Follow the Reaper", "album", "32"},
		{"Dimmu Borgir", "Enthrone Darkness Triumphant", "album", "33"},
		{"Emperor", "In the Nightside Eclipse", "album", "34"},
		{"Mayhem", "De Mysteriis Dom Sathanas", "album", "35"},
		{"Burzum", "Filosofem", "album", "36"},
		{"Darkthrone", "A Blaze in the Northern Sky", "album", "37"},
		{"Immortal", "At the Heart of Winter", "album", "38"},
		{"Satyricon", "Dark Medieval Times", "album", "39"},
		{"Gorgoroth", "Pentagram", "album", "40"},
		{"Behemoth", "The Satanist", "album", "41"},
		{"Nile", "Annihilation of the Wicked", "album", "42"},
		{"Metallica", "Garage Days Re-Revisited", "ep", "43"},
		{"Iron Maiden", "Flight of Icarus", "single", "44"},
		{"Slayer", "Seasons in the Abyss", "single", "45"},
		{"Opeth", "Damnation", "album", "46"},
		{"Opeth", "Still Life", "album", "47"},
		{"Cradle of Filth", "Dusk... and Her Embrace: The Original Sin", "album", "48"},
		{"Dimmu Borgir and the Symphony of Eternal Darkness", "Puritanical Euphoric Misanthropia", "album", "49"},
		{"The Berzerker Presents: Extreme Audio Warfare", "The Fundamental Principles of Audio Destruction", "album", "50"},
		{"Wolves in the Throne Room", "Two Hunters", "album", "51"},
		{"Agalloch", "The Mantle", "album", "52"},
	}

	releaseIDs := make(map[string]string) // map "bandName|releaseName" to ID
	for _, release := range releases {
		id := uuid.NewString()
		key := fmt.Sprintf("%s|%s", release.bandName, release.name)
		releaseIDs[key] = id
		_, err := q.CreateRelease(ctx, CreateReleaseParams{
			ID:     id,
			BandID: bandIDs[release.bandName],
			Name:   release.name,
			Type:   release.rType,
			Number: release.number,
		})
		if err != nil {
			return fmt.Errorf("failed to create release %s: %w", release.name, err)
		}
	}

	// Projects with variety of types and statuses
	projects := []struct {
		bandName    string
		releaseName string
		name        string
		pType       string
		status      string
	}{
		{"Metallica", "Kill 'Em All", "Kill 'Em All CD Reissue", "CD", "done"},
		{"Metallica", "Ride the Lightning", "Ride the Lightning Tape", "tapes", "in-progress"},
		{"Metallica", "Master of Puppets", "Master of Puppets Deluxe Box Set", "merch", "pending"},
		{"Iron Maiden", "The Number of the Beast", "Number of the Beast CD", "CD", "done"},
		{"Iron Maiden", "Powerslave", "Powerslave Tape", "tapes", "done"},
		{"Slayer", "Reign in Blood", "Reign in Blood CD", "CD", "in-progress"},
		{"Slayer", "South of Heaven", "South of Heaven Tape", "tapes", "pending"},
		{"Megadeth", "Peace Sells... but Who's Buying?", "Peace Sells CD Remaster", "CD", "done"},
		{"Megadeth", "Rust in Peace", "Rust in Peace Tape", "tapes", "in-progress"},
		{"Black Sabbath", "Paranoid", "Paranoid Vinyl Reissue", "vinyl", "done"},
		{"Black Sabbath", "Master of Reality", "Master of Reality CD", "CD", "pending"},
		{"Judas Priest", "British Steel", "British Steel Tape", "tapes", "done"},
		{"Judas Priest", "Painkiller", "Painkiller CD", "CD", "in-progress"},
		{"Pantera", "Cowboys from Hell", "Cowboys from Hell Tape", "tapes", "done"},
		{"Pantera", "Vulgar Display of Power", "Vulgar Display of Power CD", "CD", "pending"},
		{"Sepultura", "Beneath the Remains", "Beneath the Remains Tape", "tapes", "in-progress"},
		{"Sepultura", "Arise", "Arise CD", "CD", "done"},
		{"Death", "Scream Bloody Gore", "Scream Bloody Gore Tape", "tapes", "done"},
		{"Death", "Leprosy", "Leprosy CD", "CD", "in-progress"},
		{"Morbid Angel", "Altars of Madness", "Altars of Madness Tape", "tapes", "pending"},
		{"Cannibal Corpse", "Tomb of the Mutilated", "Tomb of the Mutilated CD", "CD", "done"},
		{"Opeth", "Blackwater Park", "Blackwater Park Deluxe Edition", "CD", "in-progress"},
		{"Opeth", "Ghost Reveries", "Ghost Reveries Limited Tape", "tapes", "done"},
		{"Opeth", "Damnation", "Damnation CD", "CD", "pending"},
		{"Opeth", "Still Life", "Still Life Remaster", "CD", "in-progress"},
		{"Cradle of Filth", "Dusk... and Her Embrace: The Original Sin", "Dusk and Her Embrace Original Sin Deluxe Box Set with Bonus Tracks", "merch", "done"},
		{"Dimmu Borgir and the Symphony of Eternal Darkness", "Puritanical Euphoric Misanthropia", "Puritanical Euphoric Misanthropia Limited Edition Tape", "tapes", "pending"},
		{"The Berzerker Presents: Extreme Audio Warfare", "The Fundamental Principles of Audio Destruction", "Fundamental Principles Audio Destruction CD Reissue", "CD", "in-progress"},
	}

	for _, project := range projects {
		id := uuid.NewString()
		releaseKey := fmt.Sprintf("%s|%s", project.bandName, project.releaseName)
		_, err := q.CreateProject(ctx, CreateProjectParams{
			ID:        id,
			BandID:    bandIDs[project.bandName],
			ReleaseID: releaseIDs[releaseKey],
			Name:      project.name,
			Type:      project.pType,
			Status:    project.status,
		})
		if err != nil {
			return fmt.Errorf("failed to create project %s: %w", project.name, err)
		}
	}

	return nil
}
