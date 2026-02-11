-- Down
DELETE FROM prompt_versions WHERE template_id IN (
    SELECT id FROM prompt_templates WHERE name IN (
        'ba_agent_index_system',
        'ba_agent_index_instruction',
        'ba_agent_outline_system',
        'ba_agent_outline_instruction',
        'ba_agent_full_system',
        'ba_agent_full_instruction'
    )
);

DELETE FROM prompt_templates WHERE name IN (
    'ba_agent_index_system',
    'ba_agent_index_instruction',
    'ba_agent_outline_system',
    'ba_agent_outline_instruction',
    'ba_agent_full_system',
    'ba_agent_full_instruction'
);
