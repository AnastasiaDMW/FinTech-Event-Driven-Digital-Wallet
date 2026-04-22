package com.example.balance.converter;

import com.example.balance.model.Category;
import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter(autoApply = false)
public class CategoryConverter implements AttributeConverter<Category, String> {

    @Override
    public String convertToDatabaseColumn(Category attribute) {
        return attribute == null ? null : attribute.name().toLowerCase();
    }

    @Override
    public Category convertToEntityAttribute(String dbData) {
        return dbData == null ? null : Category.valueOf(dbData.toUpperCase());
    }
}