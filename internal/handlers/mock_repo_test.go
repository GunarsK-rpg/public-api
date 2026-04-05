package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// mockRepo implements repository.Repository for handler testing.
// Each method delegates to a function field; returns "not implemented" if nil.
type mockRepo struct {
	// Classifiers
	getClassifiersFilteredFunc func(ctx context.Context, auth repository.AuthContext, filter json.RawMessage) (json.RawMessage, error)
	getSourceBooksFunc         func(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error)
	validateHeroAccessFunc     func(ctx context.Context, auth repository.AuthContext, heroID int64) error

	// Heroes - Core CRUD
	getHeroesFunc    func(ctx context.Context, auth repository.AuthContext, campaignID *int64) (json.RawMessage, error)
	getHeroFunc      func(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error)
	getHeroSheetFunc func(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error)
	upsertHeroFunc   func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	deleteHeroFunc   func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)

	// Heroes - Sub-resource getters
	getHeroAttributesFunc   func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroDefensesFunc     func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroDerivedStatsFunc func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroSkillsFunc       func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroExpertisesFunc   func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroTalentsFunc      func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroEquipmentFunc    func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroConditionsFunc   func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroInjuriesFunc     func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroGoalsFunc        func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroConnectionsFunc  func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroNotesFunc        func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getHeroCulturesFunc     func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)

	// Heroes - Sub-resource upserts
	upsertHeroAttributeFunc   func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroDefenseFunc     func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroDerivedStatFunc func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroSkillFunc       func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroExpertiseFunc   func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroTalentFunc      func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroEquipmentFunc   func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroConditionFunc   func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroInjuryFunc      func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroGoalFunc        func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroConnectionFunc  func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroNoteFunc        func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	upsertHeroCultureFunc     func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)

	// Heroes - Resource patches
	patchHeroHealthFunc      func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	patchHeroFocusFunc       func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	patchHeroInvestitureFunc func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	patchHeroCurrencyFunc    func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)

	// Heroes - Equipment modifications
	addEquipmentModificationFunc    func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	removeEquipmentModificationFunc func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)

	// Heroes - Favorite actions
	addFavoriteActionFunc    func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	removeFavoriteActionFunc func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)

	// Heroes - Avatar
	upsertHeroAvatarFunc func(ctx context.Context, auth repository.AuthContext, heroID int64, avatarKey string) (*string, error)
	deleteHeroAvatarFunc func(ctx context.Context, auth repository.AuthContext, heroID int64) (*string, error)

	// Heroes - Sub-resource deletes
	deleteHeroAttributeFunc   func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroDefenseFunc     func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroDerivedStatFunc func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroSkillFunc       func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroExpertiseFunc   func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroTalentFunc      func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroEquipmentFunc   func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroConditionFunc   func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroInjuryFunc      func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroGoalFunc        func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroConnectionFunc  func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroNoteFunc        func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	deleteHeroCultureFunc     func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)

	// Campaigns
	getCampaignsFunc             func(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error)
	getCampaignFunc              func(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error)
	getCampaignByCodeFunc        func(ctx context.Context, auth repository.AuthContext, code string) (json.RawMessage, error)
	upsertCampaignFunc           func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	deleteCampaignFunc           func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)
	removeHeroFromCampaignFunc   func(ctx context.Context, auth repository.AuthContext, heroID int64, campaignID int64) (bool, error)
	getCampaignSourceBookIDsFunc func(ctx context.Context, auth repository.AuthContext, campaignID int64) ([]int64, error)

	// Combat - NPCs
	getNpcOptionsFunc   func(ctx context.Context, auth repository.AuthContext, campaignID int64) (json.RawMessage, error)
	getNpcLibraryFunc   func(ctx context.Context, auth repository.AuthContext, campaignID int64) (json.RawMessage, error)
	getNpcFunc          func(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (json.RawMessage, error)
	getNpcByIDFunc      func(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error)
	upsertNpcFunc       func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	deleteNpcFunc       func(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (bool, error)
	upsertNpcAvatarFunc func(ctx context.Context, auth repository.AuthContext, npcID int64, campaignID int64, avatarKey string) (*string, error)
	deleteNpcAvatarFunc func(ctx context.Context, auth repository.AuthContext, npcID int64, campaignID int64) (*string, error)

	// Combat - Encounters
	getCombatsFunc     func(ctx context.Context, auth repository.AuthContext, campaignID int64) (json.RawMessage, error)
	getCombatFunc      func(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (json.RawMessage, error)
	upsertCombatFunc   func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	deleteCombatFunc   func(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (bool, error)
	endCombatRoundFunc func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)

	// Combat - NPC instances
	getNpcInstanceFunc    func(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error)
	createNpcInstanceFunc func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	patchNpcInstanceFunc  func(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error)
	deleteNpcInstanceFunc func(ctx context.Context, auth repository.AuthContext, id int64) (bool, error)

	// Combat - Companions
	getHeroNpcInstancesFunc    func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
	getCompanionNpcOptionsFunc func(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error)
}

var errNotImplemented = errors.New("not implemented")

// =============================================================================
// Classifiers
// =============================================================================

func (m *mockRepo) GetClassifiersFiltered(ctx context.Context, auth repository.AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	if m.getClassifiersFilteredFunc != nil {
		return m.getClassifiersFilteredFunc(ctx, auth, filter)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetSourceBooks(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error) {
	if m.getSourceBooksFunc != nil {
		return m.getSourceBooksFunc(ctx, auth)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) ValidateHeroAccess(ctx context.Context, auth repository.AuthContext, heroID int64) error {
	if m.validateHeroAccessFunc != nil {
		return m.validateHeroAccessFunc(ctx, auth, heroID)
	}
	return nil
}

// =============================================================================
// Heroes - Core CRUD
// =============================================================================

func (m *mockRepo) GetHeroes(ctx context.Context, auth repository.AuthContext, campaignID *int64) (json.RawMessage, error) {
	if m.getHeroesFunc != nil {
		return m.getHeroesFunc(ctx, auth, campaignID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHero(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error) {
	if m.getHeroFunc != nil {
		return m.getHeroFunc(ctx, auth, id)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroSheet(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error) {
	if m.getHeroSheetFunc != nil {
		return m.getHeroSheetFunc(ctx, auth, id)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHero(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroFunc != nil {
		return m.upsertHeroFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteHero(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroFunc != nil {
		return m.deleteHeroFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

// =============================================================================
// Heroes - Sub-resource getters
// =============================================================================

func (m *mockRepo) GetHeroAttributes(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroAttributesFunc != nil {
		return m.getHeroAttributesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroDefenses(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroDefensesFunc != nil {
		return m.getHeroDefensesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroDerivedStats(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroDerivedStatsFunc != nil {
		return m.getHeroDerivedStatsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroSkills(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroSkillsFunc != nil {
		return m.getHeroSkillsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroExpertises(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroExpertisesFunc != nil {
		return m.getHeroExpertisesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroTalents(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroTalentsFunc != nil {
		return m.getHeroTalentsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroEquipment(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroEquipmentFunc != nil {
		return m.getHeroEquipmentFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroConditions(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroConditionsFunc != nil {
		return m.getHeroConditionsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroInjuries(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroInjuriesFunc != nil {
		return m.getHeroInjuriesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroGoals(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroGoalsFunc != nil {
		return m.getHeroGoalsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroConnections(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroConnectionsFunc != nil {
		return m.getHeroConnectionsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroNotes(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroNotesFunc != nil {
		return m.getHeroNotesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetHeroCultures(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroCulturesFunc != nil {
		return m.getHeroCulturesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Heroes - Sub-resource upserts
// =============================================================================

func (m *mockRepo) UpsertHeroAttribute(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroAttributeFunc != nil {
		return m.upsertHeroAttributeFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroDefense(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroDefenseFunc != nil {
		return m.upsertHeroDefenseFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroDerivedStat(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroDerivedStatFunc != nil {
		return m.upsertHeroDerivedStatFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroSkill(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroSkillFunc != nil {
		return m.upsertHeroSkillFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroExpertise(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroExpertiseFunc != nil {
		return m.upsertHeroExpertiseFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroTalent(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroTalentFunc != nil {
		return m.upsertHeroTalentFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroEquipment(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroEquipmentFunc != nil {
		return m.upsertHeroEquipmentFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroCondition(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroConditionFunc != nil {
		return m.upsertHeroConditionFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroInjury(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroInjuryFunc != nil {
		return m.upsertHeroInjuryFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroGoal(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroGoalFunc != nil {
		return m.upsertHeroGoalFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroConnection(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroConnectionFunc != nil {
		return m.upsertHeroConnectionFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroNote(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroNoteFunc != nil {
		return m.upsertHeroNoteFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertHeroCulture(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertHeroCultureFunc != nil {
		return m.upsertHeroCultureFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Heroes - Resource patches
// =============================================================================

func (m *mockRepo) PatchHeroHealth(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.patchHeroHealthFunc != nil {
		return m.patchHeroHealthFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) PatchHeroFocus(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.patchHeroFocusFunc != nil {
		return m.patchHeroFocusFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) PatchHeroInvestiture(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.patchHeroInvestitureFunc != nil {
		return m.patchHeroInvestitureFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) PatchHeroCurrency(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.patchHeroCurrencyFunc != nil {
		return m.patchHeroCurrencyFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Heroes - Equipment modifications
// =============================================================================

func (m *mockRepo) AddEquipmentModification(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.addEquipmentModificationFunc != nil {
		return m.addEquipmentModificationFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) RemoveEquipmentModification(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.removeEquipmentModificationFunc != nil {
		return m.removeEquipmentModificationFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

// =============================================================================
// Heroes - Favorite actions
// =============================================================================

func (m *mockRepo) AddFavoriteAction(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.addFavoriteActionFunc != nil {
		return m.addFavoriteActionFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) RemoveFavoriteAction(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.removeFavoriteActionFunc != nil {
		return m.removeFavoriteActionFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

// =============================================================================
// Heroes - Avatar
// =============================================================================

func (m *mockRepo) UpsertHeroAvatar(ctx context.Context, auth repository.AuthContext, heroID int64, avatarKey string) (*string, error) {
	if m.upsertHeroAvatarFunc != nil {
		return m.upsertHeroAvatarFunc(ctx, auth, heroID, avatarKey)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteHeroAvatar(ctx context.Context, auth repository.AuthContext, heroID int64) (*string, error) {
	if m.deleteHeroAvatarFunc != nil {
		return m.deleteHeroAvatarFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Heroes - Sub-resource deletes
// =============================================================================

func (m *mockRepo) DeleteHeroAttribute(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroAttributeFunc != nil {
		return m.deleteHeroAttributeFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroDefense(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroDefenseFunc != nil {
		return m.deleteHeroDefenseFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroDerivedStat(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroDerivedStatFunc != nil {
		return m.deleteHeroDerivedStatFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroSkill(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroSkillFunc != nil {
		return m.deleteHeroSkillFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroExpertise(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroExpertiseFunc != nil {
		return m.deleteHeroExpertiseFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroTalent(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroTalentFunc != nil {
		return m.deleteHeroTalentFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroEquipment(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroEquipmentFunc != nil {
		return m.deleteHeroEquipmentFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroCondition(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroConditionFunc != nil {
		return m.deleteHeroConditionFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroInjury(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroInjuryFunc != nil {
		return m.deleteHeroInjuryFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroGoal(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroGoalFunc != nil {
		return m.deleteHeroGoalFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroConnection(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroConnectionFunc != nil {
		return m.deleteHeroConnectionFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroNote(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroNoteFunc != nil {
		return m.deleteHeroNoteFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) DeleteHeroCulture(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteHeroCultureFunc != nil {
		return m.deleteHeroCultureFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

// =============================================================================
// Campaigns
// =============================================================================

func (m *mockRepo) GetCampaigns(ctx context.Context, auth repository.AuthContext) (json.RawMessage, error) {
	if m.getCampaignsFunc != nil {
		return m.getCampaignsFunc(ctx, auth)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetCampaign(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error) {
	if m.getCampaignFunc != nil {
		return m.getCampaignFunc(ctx, auth, id)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetCampaignByCode(ctx context.Context, auth repository.AuthContext, code string) (json.RawMessage, error) {
	if m.getCampaignByCodeFunc != nil {
		return m.getCampaignByCodeFunc(ctx, auth, code)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertCampaign(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertCampaignFunc != nil {
		return m.upsertCampaignFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteCampaign(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteCampaignFunc != nil {
		return m.deleteCampaignFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

func (m *mockRepo) RemoveHeroFromCampaign(ctx context.Context, auth repository.AuthContext, heroID int64, campaignID int64) (bool, error) {
	if m.removeHeroFromCampaignFunc != nil {
		return m.removeHeroFromCampaignFunc(ctx, auth, heroID, campaignID)
	}
	return false, errNotImplemented
}

func (m *mockRepo) GetCampaignSourceBookIDs(ctx context.Context, auth repository.AuthContext, campaignID int64) ([]int64, error) {
	if m.getCampaignSourceBookIDsFunc != nil {
		return m.getCampaignSourceBookIDsFunc(ctx, auth, campaignID)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Combat - NPCs
// =============================================================================

func (m *mockRepo) GetNpcOptions(ctx context.Context, auth repository.AuthContext, campaignID int64) (json.RawMessage, error) {
	if m.getNpcOptionsFunc != nil {
		return m.getNpcOptionsFunc(ctx, auth, campaignID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetNpcLibrary(ctx context.Context, auth repository.AuthContext, campaignID int64) (json.RawMessage, error) {
	if m.getNpcLibraryFunc != nil {
		return m.getNpcLibraryFunc(ctx, auth, campaignID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetNpc(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (json.RawMessage, error) {
	if m.getNpcFunc != nil {
		return m.getNpcFunc(ctx, auth, id, campaignID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetNpcByID(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error) {
	if m.getNpcByIDFunc != nil {
		return m.getNpcByIDFunc(ctx, auth, id)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertNpc(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertNpcFunc != nil {
		return m.upsertNpcFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteNpc(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (bool, error) {
	if m.deleteNpcFunc != nil {
		return m.deleteNpcFunc(ctx, auth, id, campaignID)
	}
	return false, errNotImplemented
}

func (m *mockRepo) UpsertNpcAvatar(ctx context.Context, auth repository.AuthContext, npcID int64, campaignID int64, avatarKey string) (*string, error) {
	if m.upsertNpcAvatarFunc != nil {
		return m.upsertNpcAvatarFunc(ctx, auth, npcID, campaignID, avatarKey)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteNpcAvatar(ctx context.Context, auth repository.AuthContext, npcID int64, campaignID int64) (*string, error) {
	if m.deleteNpcAvatarFunc != nil {
		return m.deleteNpcAvatarFunc(ctx, auth, npcID, campaignID)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Combat - Encounters
// =============================================================================

func (m *mockRepo) GetCombats(ctx context.Context, auth repository.AuthContext, campaignID int64) (json.RawMessage, error) {
	if m.getCombatsFunc != nil {
		return m.getCombatsFunc(ctx, auth, campaignID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetCombat(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (json.RawMessage, error) {
	if m.getCombatFunc != nil {
		return m.getCombatFunc(ctx, auth, id, campaignID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) UpsertCombat(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.upsertCombatFunc != nil {
		return m.upsertCombatFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteCombat(ctx context.Context, auth repository.AuthContext, id int64, campaignID int64) (bool, error) {
	if m.deleteCombatFunc != nil {
		return m.deleteCombatFunc(ctx, auth, id, campaignID)
	}
	return false, errNotImplemented
}

func (m *mockRepo) EndCombatRound(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.endCombatRoundFunc != nil {
		return m.endCombatRoundFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

// =============================================================================
// Combat - NPC instances
// =============================================================================

func (m *mockRepo) GetNpcInstance(ctx context.Context, auth repository.AuthContext, id int64) (json.RawMessage, error) {
	if m.getNpcInstanceFunc != nil {
		return m.getNpcInstanceFunc(ctx, auth, id)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) CreateNpcInstance(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.createNpcInstanceFunc != nil {
		return m.createNpcInstanceFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) PatchNpcInstance(ctx context.Context, auth repository.AuthContext, data json.RawMessage) (json.RawMessage, error) {
	if m.patchNpcInstanceFunc != nil {
		return m.patchNpcInstanceFunc(ctx, auth, data)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) DeleteNpcInstance(ctx context.Context, auth repository.AuthContext, id int64) (bool, error) {
	if m.deleteNpcInstanceFunc != nil {
		return m.deleteNpcInstanceFunc(ctx, auth, id)
	}
	return false, errNotImplemented
}

// =============================================================================
// Combat - Companions
// =============================================================================

func (m *mockRepo) GetHeroNpcInstances(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getHeroNpcInstancesFunc != nil {
		return m.getHeroNpcInstancesFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}

func (m *mockRepo) GetCompanionNpcOptions(ctx context.Context, auth repository.AuthContext, heroID int64) (json.RawMessage, error) {
	if m.getCompanionNpcOptionsFunc != nil {
		return m.getCompanionNpcOptionsFunc(ctx, auth, heroID)
	}
	return nil, errNotImplemented
}
