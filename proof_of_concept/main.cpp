#include <SFML/Graphics.hpp>
#include <optional>

// pimpalysas

int main()
{
    // Create window
    sf::RenderWindow window(sf::VideoMode({800, 600}), "SFML Background + Button");

    // Load background texture
    sf::Texture bgTexture;
    if (!bgTexture.loadFromFile("cute_image.jpg")) {
        return -1; // fail if missing
    }

    // Sprite for background
    sf::Sprite bgSprite(bgTexture);

    // Scale the sprite to fill the window
    sf::Vector2u texSize = bgTexture.getSize();
    sf::Vector2u winSize = window.getSize();

    float scaleX = float(winSize.x) / texSize.x;
    float scaleY = float(winSize.y) / texSize.y;
    bgSprite.setScale(sf::Vector2f(scaleX, scaleY));

    // Button
    sf::RectangleShape button({200.f, 60.f});
    button.setPosition({(winSize.x - 200.f)/2, (winSize.y - 60.f)/2});
    button.setFillColor(sf::Color::Blue);

    // Load font
    sf::Font font("arial.ttf");

    // Button text
    sf::Text text(font, "Click me", 24);
    text.setFillColor(sf::Color::White);
    text.setPosition(sf::Vector2f(
        button.getPosition().x + 40.f,
        button.getPosition().y + 15.f
        )
    );

    bool isRed = false;

    while (window.isOpen())
    {
        while (const std::optional event = window.pollEvent())
        {
            if (event->is<sf::Event::Closed>())
                window.close();

            if (const auto* mouse = event->getIf<sf::Event::MouseButtonPressed>())
            {
                if (mouse->button == sf::Mouse::Button::Left)
                {
                    sf::Vector2f pos = {(float)mouse->position.x, (float)mouse->position.y};
                    if (button.getGlobalBounds().contains(pos))
                    {
                        isRed = !isRed;
                        button.setFillColor(isRed ? sf::Color::Red : sf::Color::Blue);
                    }
                }
            }
        }

        window.clear();
        window.draw(bgSprite);  // draw background first
        window.draw(button);
        window.draw(text);
        window.display();
    }
}
