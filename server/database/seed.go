package database

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/vague2k/blkhell/common"
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
	// Use random countries from our data package for variety
	bandNames := []string{
		"Metallica",
		"Iron Maiden",
		"Slayer",
		"Opeth",
		"Death",
		"Black Sabbath",
		"Pantera",
		"Sepultura",
		"Dimmu Borgir",
		"Cradle of Filth",
		"Wolves in the Throne Room",
		"The Berzerker Presents: Extreme Audio Warfare",
	}

	bands := make([]struct {
		name    string
		country string
	}, len(bandNames))

	countryCodes := make([]string, 0, len(common.Countries))
	for code := range common.Countries {
		countryCodes = append(countryCodes, code)
	}
	for i, name := range bandNames {
		// Assign a random country from our data package
		randCode := countryCodes[rand.Intn(len(countryCodes))]
		bands[i] = struct {
			name    string
			country string
		}{
			name:    name,
			country: randCode,
		}
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
		{"Metallica", "Kill 'Em All", "album", "1"},
		{"Metallica", "Ride the Lightning", "album", "2"},
		{"Metallica", "Master of Puppets", "album", "3"},
		{"Metallica", "...And Justice for All", "album", "4"},
		{"Metallica", "Metallica (Black Album)", "album", "5"},
		{"Metallica", "Load", "album", "6"},
		{"Metallica", "Garage Days Re-Revisited", "ep", "7"},
		{"Metallica", "S&M Live", "compilation", "8"},
		{"Iron Maiden", "Iron Maiden", "album", "9"},
		{"Iron Maiden", "Killers", "album", "10"},
		{"Iron Maiden", "The Number of the Beast", "album", "11"},
		{"Iron Maiden", "Piece of Mind", "album", "12"},
		{"Iron Maiden", "Powerslave", "album", "13"},
		{"Iron Maiden", "Somewhere in Time", "album", "14"},
		{"Iron Maiden", "Seventh Son of a Seventh Son", "album", "15"},
		{"Iron Maiden", "Fear of the Dark", "album", "16"},
		{"Iron Maiden", "Flight of Icarus", "single", "17"},
		{"Iron Maiden", "Live After Death", "compilation", "18"},
		{"Slayer", "Show No Mercy", "album", "19"},
		{"Slayer", "Hell Awaits", "album", "20"},
		{"Slayer", "Reign in Blood", "album", "21"},
		{"Slayer", "South of Heaven", "album", "22"},
		{"Slayer", "Seasons in the Abyss", "album", "23"},
		{"Slayer", "Divine Intervention", "album", "24"},
		{"Slayer", "Haunting the Chapel", "ep", "25"},
		{"Opeth", "Orchid", "album", "26"},
		{"Opeth", "Morningrise", "album", "27"},
		{"Opeth", "My Arms, Your Hearse", "album", "28"},
		{"Opeth", "Still Life", "album", "29"},
		{"Opeth", "Blackwater Park", "album", "30"},
		{"Opeth", "Deliverance", "album", "31"},
		{"Opeth", "Damnation", "album", "32"},
		{"Opeth", "Ghost Reveries", "album", "33"},
		{"Opeth", "Watershed", "album", "34"},
		{"Opeth", "Heritage", "album", "35"},
		{"Opeth", "Pale Communion", "album", "36"},
		{"Opeth", "In Cauda Venenum", "album", "37"},
		{"Death", "Scream Bloody Gore", "album", "38"},
		{"Death", "Leprosy", "album", "39"},
		{"Death", "Spiritual Healing", "album", "40"},
		{"Death", "Human", "album", "41"},
		{"Death", "Individual Thought Patterns", "album", "42"},
		{"Death", "Symbolic", "album", "43"},
		{"Death", "The Sound of Perseverance", "album", "44"},
		{"Black Sabbath", "Black Sabbath", "album", "45"},
		{"Black Sabbath", "Paranoid", "album", "46"},
		{"Black Sabbath", "Master of Reality", "album", "47"},
		{"Black Sabbath", "Vol. 4", "album", "48"},
		{"Black Sabbath", "Sabbath Bloody Sabbath", "album", "49"},
		{"Black Sabbath", "Sabotage", "album", "50"},
		{"Black Sabbath", "Heaven and Hell", "album", "51"},
		{"Black Sabbath", "The Eternal Idol", "album", "52"},
		{"Pantera", "Cowboys from Hell", "album", "53"},
		{"Pantera", "Vulgar Display of Power", "album", "54"},
		{"Pantera", "Far Beyond Driven", "album", "55"},
		{"Pantera", "The Great Southern Trendkill", "album", "56"},
		{"Pantera", "Reinventing the Steel", "album", "57"},
		{"Pantera", "Official Live: 101 Proof", "compilation", "58"},
		{"Sepultura", "Morbid Visions", "album", "59"},
		{"Sepultura", "Schizophrenia", "album", "60"},
		{"Sepultura", "Beneath the Remains", "album", "61"},
		{"Sepultura", "Arise", "album", "62"},
		{"Sepultura", "Chaos A.D.", "album", "63"},
		{"Sepultura", "Roots", "album", "64"},
		{"Sepultura", "Against", "album", "65"},
		{"Dimmu Borgir", "For All Tid", "album", "66"},
		{"Dimmu Borgir", "Stormblåst", "album", "67"},
		{"Dimmu Borgir", "Enthrone Darkness Triumphant", "album", "68"},
		{"Dimmu Borgir", "Spiritual Black Dimensions", "album", "69"},
		{"Dimmu Borgir", "Puritanical Euphoric Misanthropia", "album", "70"},
		{"Dimmu Borgir", "Death Cult Armageddon", "album", "71"},
		{"Cradle of Filth", "The Principle of Evil Made Flesh", "album", "72"},
		{"Cradle of Filth", "Vempire", "ep", "73"},
		{"Cradle of Filth", "Dusk and Her Embrace", "album", "74"},
		{"Cradle of Filth", "Cruelty and the Beast", "album", "75"},
		{"Cradle of Filth", "Midian", "album", "76"},
		{"Wolves in the Throne Room", "Diadem of 12 Stars", "album", "77"},
		{"Wolves in the Throne Room", "Two Hunters", "album", "78"},
		{"Wolves in the Throne Room", "Black Cascade", "album", "79"},
		{"Wolves in the Throne Room", "Celestial Lineage", "album", "80"},
		{"The Berzerker Presents: Extreme Audio Warfare", "The Fundamental Principles of Audio Destruction", "album", "81"},
		{"The Berzerker Presents: Extreme Audio Warfare", "Dissimulate", "album", "82"},
		{"The Berzerker Presents: Extreme Audio Warfare", "World of Lies", "album", "83"},
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

	// Projects with variety of types and statuses - multiple per release
	projects := []struct {
		bandName    string
		releaseName string
		name        string
		pType       string
		status      string
	}{
		{"Metallica", "Kill 'Em All", "Kill 'Em All CD Reissue", "CD", "done"},
		{"Metallica", "Kill 'Em All", "Kill 'Em All Tape", "tapes", "done"},
		{"Metallica", "Ride the Lightning", "Ride the Lightning CD", "CD", "in-progress"},
		{"Metallica", "Ride the Lightning", "Ride the Lightning Vinyl", "vinyl", "pending"},
		{"Metallica", "Master of Puppets", "Master of Puppets Deluxe Box Set", "merch", "done"},
		{"Metallica", "Master of Puppets", "Master of Puppets CD", "CD", "done"},
		{"Metallica", "...And Justice for All", "Justice CD", "CD", "in-progress"},
		{"Metallica", "Metallica (Black Album)", "Black Album Vinyl Reissue", "vinyl", "pending"},
		{"Metallica", "Load", "Load Tape", "tapes", "pending"},
		{"Iron Maiden", "The Number of the Beast", "Number of the Beast CD", "CD", "done"},
		{"Iron Maiden", "The Number of the Beast", "Number of the Beast Vinyl", "vinyl", "done"},
		{"Iron Maiden", "Powerslave", "Powerslave Tape", "tapes", "done"},
		{"Iron Maiden", "Powerslave", "Powerslave CD Remaster", "CD", "in-progress"},
		{"Iron Maiden", "Somewhere in Time", "Somewhere in Time CD", "CD", "pending"},
		{"Iron Maiden", "Seventh Son of a Seventh Son", "Seventh Son Vinyl", "vinyl", "done"},
		{"Iron Maiden", "Fear of the Dark", "Fear of the Dark CD", "CD", "in-progress"},
		{"Iron Maiden", "Live After Death", "Live After Death Box Set", "merch", "done"},
		{"Slayer", "Reign in Blood", "Reign in Blood CD", "CD", "done"},
		{"Slayer", "Reign in Blood", "Reign in Blood Tape", "tapes", "done"},
		{"Slayer", "Reign in Blood", "Reign in Blood Vinyl Reissue", "vinyl", "in-progress"},
		{"Slayer", "South of Heaven", "South of Heaven CD", "CD", "pending"},
		{"Slayer", "South of Heaven", "South of Heaven Tape", "tapes", "done"},
		{"Slayer", "Seasons in the Abyss", "Seasons CD", "CD", "in-progress"},
		{"Slayer", "Divine Intervention", "Divine Intervention Tape", "tapes", "pending"},
		{"Opeth", "Blackwater Park", "Blackwater Park Deluxe Edition", "CD", "done"},
		{"Opeth", "Blackwater Park", "Blackwater Park Vinyl", "vinyl", "done"},
		{"Opeth", "Blackwater Park", "Blackwater Park Tape", "tapes", "in-progress"},
		{"Opeth", "Ghost Reveries", "Ghost Reveries Limited Tape", "tapes", "done"},
		{"Opeth", "Ghost Reveries", "Ghost Reveries CD", "CD", "done"},
		{"Opeth", "Damnation", "Damnation CD", "CD", "pending"},
		{"Opeth", "Damnation", "Damnation Vinyl", "vinyl", "in-progress"},
		{"Opeth", "Still Life", "Still Life Remaster CD", "CD", "done"},
		{"Opeth", "Still Life", "Still Life Vinyl", "vinyl", "pending"},
		{"Opeth", "Deliverance", "Deliverance CD", "CD", "in-progress"},
		{"Opeth", "Watershed", "Watershed Deluxe", "merch", "done"},
		{"Opeth", "Heritage", "Heritage Vinyl", "vinyl", "pending"},
		{"Opeth", "Pale Communion", "Pale Communion CD", "CD", "in-progress"},
		{"Opeth", "In Cauda Venenum", "In Cauda Venenum Box Set", "merch", "pending"},
		{"Death", "Scream Bloody Gore", "Scream Bloody Gore Tape", "tapes", "done"},
		{"Death", "Scream Bloody Gore", "Scream Bloody Gore CD Reissue", "CD", "done"},
		{"Death", "Leprosy", "Leprosy CD", "CD", "in-progress"},
		{"Death", "Leprosy", "Leprosy Vinyl", "vinyl", "pending"},
		{"Death", "Human", "Human CD", "CD", "done"},
		{"Death", "Symbolic", "Symbolic Tape", "tapes", "in-progress"},
		{"Death", "The Sound of Perseverance", "Sound of Perseverance CD", "CD", "pending"},
		{"Black Sabbath", "Paranoid", "Paranoid Vinyl Reissue", "vinyl", "done"},
		{"Black Sabbath", "Paranoid", "Paranoid CD", "CD", "done"},
		{"Black Sabbath", "Master of Reality", "Master of Reality CD", "CD", "in-progress"},
		{"Black Sabbath", "Master of Reality", "Master of Reality Tape", "tapes", "pending"},
		{"Black Sabbath", "Vol. 4", "Vol. 4 Vinyl", "vinyl", "done"},
		{"Black Sabbath", "Sabbath Bloody Sabbath", "Sabbath Bloody Sabbath CD", "CD", "in-progress"},
		{"Black Sabbath", "Heaven and Hell", "Heaven and Hell Deluxe", "merch", "pending"},
		{"Pantera", "Cowboys from Hell", "Cowboys from Hell Tape", "tapes", "done"},
		{"Pantera", "Cowboys from Hell", "Cowboys from Hell CD", "CD", "done"},
		{"Pantera", "Vulgar Display of Power", "Vulgar Display CD", "CD", "in-progress"},
		{"Pantera", "Vulgar Display of Power", "Vulgar Display Vinyl", "vinyl", "pending"},
		{"Pantera", "Far Beyond Driven", "Far Beyond Driven CD", "CD", "done"},
		{"Pantera", "The Great Southern Trendkill", "Trendkill Tape", "tapes", "in-progress"},
		{"Sepultura", "Beneath the Remains", "Beneath the Remains Tape", "tapes", "done"},
		{"Sepultura", "Beneath the Remains", "Beneath the Remains CD", "CD", "in-progress"},
		{"Sepultura", "Arise", "Arise CD", "CD", "done"},
		{"Sepultura", "Arise", "Arise Vinyl", "vinyl", "pending"},
		{"Sepultura", "Chaos A.D.", "Chaos A.D. CD", "CD", "in-progress"},
		{"Sepultura", "Roots", "Roots Tape", "tapes", "done"},
		{"Sepultura", "Roots", "Roots Vinyl Reissue", "vinyl", "pending"},
		{"Dimmu Borgir", "Enthrone Darkness Triumphant", "Enthrone Darkness CD", "CD", "done"},
		{"Dimmu Borgir", "Enthrone Darkness Triumphant", "Enthrone Darkness Tape", "tapes", "in-progress"},
		{"Dimmu Borgir", "Puritanical Euphoric Misanthropia", "Puritanical CD", "CD", "pending"},
		{"Dimmu Borgir", "Death Cult Armageddon", "Death Cult Vinyl", "vinyl", "done"},
		{"Dimmu Borgir", "Spiritual Black Dimensions", "Spiritual Black Tape", "tapes", "in-progress"},
		{"Cradle of Filth", "Dusk and Her Embrace", "Dusk and Her Embrace Deluxe Box Set", "merch", "done"},
		{"Cradle of Filth", "Dusk and Her Embrace", "Dusk and Her Embrace CD", "CD", "done"},
		{"Cradle of Filth", "Cruelty and the Beast", "Cruelty CD", "CD", "in-progress"},
		{"Cradle of Filth", "Midian", "Midian Vinyl", "vinyl", "pending"},
		{"Wolves in the Throne Room", "Two Hunters", "Two Hunters Vinyl", "vinyl", "done"},
		{"Wolves in the Throne Room", "Two Hunters", "Two Hunters CD", "CD", "in-progress"},
		{"Wolves in the Throne Room", "Black Cascade", "Black Cascade Tape", "tapes", "pending"},
		{"Wolves in the Throne Room", "Celestial Lineage", "Celestial Lineage CD", "CD", "done"},
		{"The Berzerker Presents: Extreme Audio Warfare", "The Fundamental Principles of Audio Destruction", "Fundamental Principles CD Reissue", "CD", "done"},
		{"The Berzerker Presents: Extreme Audio Warfare", "The Fundamental Principles of Audio Destruction", "Fundamental Principles Tape Limited Edition", "tapes", "in-progress"},
		{"The Berzerker Presents: Extreme Audio Warfare", "Dissimulate", "Dissimulate CD", "CD", "pending"},
		{"The Berzerker Presents: Extreme Audio Warfare", "World of Lies", "World of Lies Vinyl", "vinyl", "done"},
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
